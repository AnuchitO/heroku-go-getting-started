package profile

import (
	"context"
	"reflect"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type storage struct {
	db *mongo.Database
}

func NewStorage(db *mongo.Database) *storage {
	return &storage{
		db: db,
	}
}

const userCollection = "users"

func (s *storage) List() ([]Profile, error) {
	query := bson.M{}
	result, err := s.db.Collection(userCollection).Find(context.TODO(), query, options.Find())
	if err != nil {
		return nil, err
	}

	var profiles []Profile
	err = result.All(context.TODO(), &profiles)
	return profiles, err
}

func (s *storage) AboutMeUpdate(id string, about aboutme) error {
	update := bson.M{"$set": bson.M{
		"about_me":      about.AboutMe,
		"social_medias": about.SocialMedia,
		"tags":          about.Tags,
		"updated_at":    time.Now(),
	}}

	result, err := s.db.Collection(userCollection).UpdateByID(context.TODO(), id, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount <= 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (s *storage) UpdateProfileSkills(id string, set skillset, skills []Skill) error {
	var err error
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{string(set): skills}}
	_, err = s.db.Collection(userCollection).UpdateOne(context.TODO(), filter, update)

	return err
}

func (s *storage) GetByID(ctx context.Context, id string) (*Profile, error) {
	filter := bson.M{"_id": id}
	var profile Profile
	if err := s.db.Collection(userCollection).FindOne(ctx, filter, options.FindOne()).Decode(&profile); err != nil {
		return nil, err
	}
	if reflect.DeepEqual(profile, Profile{}) {
		return nil, mongo.ErrNoDocuments
	}
	return &profile, nil
}

func (s *storage) GetSkills(ctx context.Context, id string, kind string) (*SkillsByUser, error) {
	if len(kind) < 1 {
		return &SkillsByUser{}, mongo.WriteError{Code: 400, Message: "Missing Skill Kind"}
	}
	var pipeline []primitive.M

	if kind != "technical" && kind != "soft" {
		return nil, ErrInvalidKindOfSkill
	}

	if kind == "technical" {
		pipeline = Pipeliner(id, TechnicalSkillQuerySet())
	}

	if kind == "soft" {
		pipeline = Pipeliner(id, SoftSkillQuerySet())
	}

	cursor, err := s.db.Collection(userCollection).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var results []SkillsByUser
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	skills := results[0].Skills

	sort.SliceStable(skills, func(i, j int) bool {
		return skills[i].Score > skills[j].Score
	})

	res := results[0]
	res.Skills = skills

	return &res, nil
}

func Pipeliner(oid string, set SkillQuerySet) []bson.M {
	matchStage := MatchStage(oid)
	lookupStage := LookupStage(set)
	projectStage := ProjectStage(set)

	return []bson.M{matchStage, lookupStage, projectStage}
}

func TechnicalSkillQuerySet() SkillQuerySet {
	return SkillQuerySet{
		localField:    "technical_skills.skillID",
		mapInput:      "$technical_skills",
		mapInputAlias: "technical",
		refSkillID:    "$$technical.skillID",
		refSkillScore: "$$technical.score",
	}
}

func SoftSkillQuerySet() SkillQuerySet {
	return SkillQuerySet{
		localField:    "soft_skills.skillID",
		mapInput:      "$soft_skills",
		mapInputAlias: "soft",
		refSkillID:    "$$soft.skillID",
		refSkillScore: "$$soft.score",
	}
}

func MatchStage(oid string) bson.M {
	return bson.M{
		"$match": bson.M{
			"_id": oid,
		},
	}
}

func LookupStage(set SkillQuerySet) bson.M {
	return bson.M{
		"$lookup": bson.M{
			"from":         "skills",
			"localField":   set.localField,
			"foreignField": "_id",
			"as":           "SkillData",
		},
	}
}

func ProjectStage(set SkillQuerySet) bson.M {
	return bson.M{
		"$project": bson.M{
			"Skills": bson.M{
				"$map": bson.M{
					"input": set.mapInput,
					"as":    set.mapInputAlias,
					"in": bson.M{
						"skill": bson.M{
							"$arrayElemAt": []interface{}{
								bson.M{
									"$filter": bson.M{
										"input": "$SkillData",
										"as":    "info",
										"cond": bson.M{
											"$eq": []interface{}{
												"$$info._id",
												set.refSkillID,
											},
										},
									},
								},
								0,
							},
						},
						"score": set.refSkillScore,
					},
				},
			},
		},
	}
}

func (s *storage) GetAll(ctx context.Context) ([]User, error) {
	query := bson.M{}
	result, err := s.db.Collection(userCollection).Find(ctx, query, options.Find())
	if err != nil {
		return nil, err
	}

	var results []User
	if err = result.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

type UserStorageError struct {
	message string
}

func (e UserStorageError) Error() string {
	return e.message
}

var UserNotFoundError = UserStorageError{message: "user not found"}

func (s *storage) GetOneByEmail(ctx context.Context, email string) (*User, error) {
	filter := bson.M{"email": email}
	var user User

	err := s.db.Collection(userCollection).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, UserNotFoundError
	}

	return &user, nil
}
