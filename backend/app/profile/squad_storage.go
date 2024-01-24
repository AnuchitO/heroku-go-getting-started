package profile

import (
	"context"
	"errors"

	"gitdev.devops.krungthai.com/aster/ariskill/app/squad"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const squadCollection = "squads"

func (s *storage) GetOneByID(ctx context.Context, id string) (*squad.Squad, error) {
	obId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var squad *squad.Squad
	err = s.db.Collection(squadCollection).FindOne(ctx, bson.M{"_id": obId}, options.FindOne()).Decode(&squad)
	if err != nil {
		return nil, err
	}

	return squad, nil
}

func (s *storage) UpdateByID(ctx context.Context, id string, updateSquad *squad.Squad) (*squad.Squad, error) {
	update := bson.M{"$set": updateSquad}
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	res, errRes := s.db.Collection(squadCollection).UpdateOne(ctx, bson.M{"_id": objId}, update)
	if errRes != nil {
		return nil, errRes
	}

	if res.MatchedCount <= 0 {
		return nil, errors.New("squad not found")
	}

	var result *squad.Squad
	err1 := s.db.Collection(squadCollection).FindOne(ctx, bson.D{{Key: "_id", Value: objId}}).Decode(&result)
	if err1 != nil {
		return nil, err1
	}

	return result, nil
}
