package skill

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type storage struct {
	db *mongo.Database
}

func NewStorage(db *mongo.Database) *storage {
	return &storage{
		db: db,
	}
}

const skillCollection = "skills"
const hardSkillCollection = "hard_skills"

func (s *storage) GetByKind(ctx context.Context, kind string) ([]Skill, error) {
	query := bson.M{}
	if kind != "" {
		query = bson.M{"kind": kind}
	}
	var sks []Skill
	cursor, err := s.db.Collection(skillCollection).Find(ctx, query)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &sks)
	return sks, err
}

func (s *storage) GetByID(ctx context.Context, id string) (Skill, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Skill{}, err
	}
	filter := bson.M{"_id": oid}

	var result Skill
	err = s.db.Collection(skillCollection).FindOne(context.Background(), filter).Decode(&result)

	return result, err
}

func (s *storage) GetByRole(ctx context.Context, role string) ([]HardSkill, error) {
	filter := bson.M{
		"jobRole": bson.M{
			"$in": []string{role},
		},
	}

	var result = []HardSkill{}
	cur, err := s.db.Collection(hardSkillCollection).Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	err = cur.All(context.Background(), &result)

	return result, err
}
