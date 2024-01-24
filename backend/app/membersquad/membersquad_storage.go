package membersquad

import (
	"context"
	"reflect"

	"gitdev.devops.krungthai.com/aster/ariskill/app/user"
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

func (s *storage) GetByID(ctx context.Context, id string) (*Member, error) {
	filter := bson.M{"_id": id}
	var entity Member

	err := s.db.Collection(userCollection).FindOne(ctx, filter, options.FindOne()).Decode(&entity)
	if err != nil {
		return nil, err
	}
	if reflect.DeepEqual(entity, user.User{}) {
		return nil, mongo.ErrNoDocuments
	}
	return &entity, nil
}

func (s *storage) UpdateMySquad(ctx context.Context, id string, update []Squad) error {
	res, err := s.db.Collection(userCollection).UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"my_squad": update}})
	if err != nil {
		return err
	}

	if res.MatchedCount < 1 {
		return mongo.WriteError{Code: 404, Message: "mongo: no documents in result"}
	}

	if res.ModifiedCount < 0 {
		return mongo.WriteError{Code: 404, Message: "mongo: document found but doesn't have any updated"}
	}
	return nil
}

func (s *storage) GetAllBySquadId(ctx context.Context, squadId primitive.ObjectID) ([]Member, error) {
	filter := bson.M{"my_squad": bson.M{"$elemMatch": bson.M{"sqid": squadId}}}

	res, err := s.db.Collection(userCollection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var users []Member
	if err := res.All(ctx, &users); err != nil {
		return nil, err
	}

	if len(users) < 1 {
		return nil, mongo.ErrNoDocuments
	}

	return users, nil
}
