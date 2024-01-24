package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const userCollection = "users"

type storage struct {
	db *mongo.Database
}

func NewStorage(db *mongo.Database) *storage {
	return &storage{
		db: db,
	}
}

type UserStorageError struct {
	message string
}

func (e UserStorageError) Error() string {
	return e.message
}

// var invalidIdError = UserStorageError{message: "invalid user id"}
var UserNotFoundError = UserStorageError{message: "user not found"}
var teamNotFoundError = UserStorageError{message: "team not found"}
var IsLeadFoundError = UserStorageError{message: "teamlead"}

// var dbConnectNotFound = CycleStorageError{message: "cannot connect to mongodb"}

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

func (s *storage) GetAllBySquadId(ctx context.Context, squadId primitive.ObjectID) ([]User, error) {
	filter := bson.M{"my_squad": bson.M{"$elemMatch": bson.M{"sqid": squadId}}}

	res, err := s.db.Collection(userCollection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var users []User
	if err := res.All(ctx, &users); err != nil {
		return nil, err
	}

	if len(users) < 1 {
		return nil, mongo.ErrNoDocuments
	}

	return users, nil
}

// Get user with user_id
func (s *storage) GetOneById(ctx context.Context, id string) (*User, error) {
	filter := bson.M{"_id": id}
	var user User

	err := s.db.Collection(userCollection).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, UserNotFoundError
	}

	return &user, nil
}

// Get user with user_id
func (s *storage) GetOneByEmail(ctx context.Context, email string) (*User, error) {
	filter := bson.M{"email": email}
	var user User

	err := s.db.Collection(userCollection).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, UserNotFoundError
	}

	return &user, nil
}

// Get HardSkill with user_id
func (s *storage) GetHardSkillById(ctx context.Context, name string) ([]User, error) {
	// Case-insensitive regex for the name
	filter := bson.M{"hardSkills": primitive.Regex{Pattern: name, Options: "i"}}

	res, err := s.db.Collection(userCollection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var users []User
	if err := res.All(ctx, &users); err != nil {
		return nil, err
	}

	if len(users) < 1 {
		return nil, mongo.ErrNoDocuments
	}

	return users, nil
}
