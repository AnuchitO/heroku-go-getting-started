package squad

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type storage struct {
	db *mongo.Database
}

func NewSquadStorage(db *mongo.Database) *storage {
	return &storage{
		db: db,
	}
}

const squadCollection = "squads"

type SquadStorageError struct {
	message string
}

func (e SquadStorageError) Error() string {
	return e.message
}

var invalidIdError = SquadStorageError{message: "invalid squad id"}
var squadNotFoundError = SquadStorageError{message: "squad not found"}

func convertIdToObjectId(id string) (*primitive.ObjectID, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return &objectId, nil
}

func (s *storage) InsertOne(profileId string, sq Squad) (*Squad, error) {
	// The insert user id is the owner
	// The insert squad skill must already has 1 rating for each skill
	// Set the rating of each skill user id to creator id

	// Loop over each skills
	for i, v := range sq.SkillsRatings {
		// Loop over each rating
		for j := range v.Ratings {
			// Set the user id
			sq.SkillsRatings[i].Ratings[j].UserId = GoogleUserId(profileId)
		}
	}
	sq.CreatedAt = time.Now()

	res, err := s.db.Collection(squadCollection).InsertOne(context.TODO(), sq, options.InsertOne())
	if err != nil {
		return nil, err
	}

	sq.Id = res.InsertedID.(primitive.ObjectID)

	return &sq, nil
}

func (s *storage) GetAll() ([]*Squad, error) {
	var squads []*Squad
	cursor, err := s.db.Collection(squadCollection).Find(context.Background(), bson.M{}, options.Find())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, squadNotFoundError
		}
		return nil, err
	}
	err = cursor.All(context.Background(), &squads)
	if err != nil {
		return nil, err
	}
	return squads, nil
}

func (s *storage) GetByFilter(filter regexFilter) ([]*Squad, error) {
	var squads []*Squad
	cursor, err := s.db.Collection(squadCollection).Find(context.Background(), filter, options.Find())
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &squads)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, squadNotFoundError
		}
		return nil, err
	}
	return squads, nil
}

func (s *storage) GetOneByID(id string) (*Squad, error) {
	objectId, err := convertIdToObjectId(id)
	if err != nil {
		return nil, invalidIdError
	}

	var squad *Squad
	err = s.db.Collection(squadCollection).FindOne(context.TODO(), bson.M{"_id": *objectId}, options.FindOne()).Decode(&squad)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, squadNotFoundError
		}
		return nil, err
	}
	return squad, nil
}

func (s *storage) UpdateOneByID(id string, updatedSquad Squad) (*Squad, error) {
	objectId, err := convertIdToObjectId(id)
	if err != nil {
		return nil, invalidIdError
	}

	update := bson.M{"$set": updatedSquad}

	result, err := s.db.Collection(squadCollection).UpdateOne(context.TODO(), bson.M{"_id": *objectId}, update, options.Update())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, squadNotFoundError
		}
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, squadNotFoundError
	}

	updatedSquad.Id = *objectId
	return &updatedSquad, nil
}

func (s *storage) DeleteByID(id string) error {
	objectId, err := convertIdToObjectId(id)
	if err != nil {
		return err
	}

	_, err = s.db.Collection(squadCollection).DeleteOne(context.TODO(), bson.M{"_id": *objectId}, options.Delete())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return squadNotFoundError
		}
		return err
	}

	return nil
}
