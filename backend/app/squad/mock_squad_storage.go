package squad

import (
	"fmt"
	"testing"

	"gitdev.devops.krungthai.com/aster/ariskill/app/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockSquadStorage struct {
	squad         []*Squad
	methodsToCall map[string]bool
	err           error
	users         []user.User
}

type mockingObjectId struct {
	hexId    string
	objectId primitive.ObjectID
}

func mockObjectId(id int) mockingObjectId {
	hexId := fmt.Sprintf("%024d", id)
	objectId, _ := primitive.ObjectIDFromHex(hexId)
	return mockingObjectId{hexId: hexId, objectId: objectId}
}

func mockGoogleUserId(id int) GoogleUserId {
	stringId := fmt.Sprintf("%021d", id)
	return GoogleUserId(stringId)
}

func (ms *mockSquadStorage) InsertOne(userId string, squadToInsert Squad) (*Squad, error) {
	ms.methodsToCall["InsertOneByID"] = true
	if ms.err != nil {
		return nil, ms.err
	}
	insertedSquad := Squad{
		Id:            ms.squad[0].Id,
		Name:          squadToInsert.Name,
		TeamleadMail:  "john.d@arise.tech",
		Description:   squadToInsert.Description,
		SkillsRatings: squadToInsert.SkillsRatings,
	}
	ms.squad = append(ms.squad, &insertedSquad)
	return ms.squad[len(ms.squad)-1], nil
}

func (ms *mockSquadStorage) GetAll() ([]*Squad, error) {
	ms.methodsToCall["GetAll"] = true
	if ms.err != nil {
		return nil, ms.err
	}

	return ms.squad, nil
}

func (ms *mockSquadStorage) GetByFilter(filter regexFilter) ([]*Squad, error) {
	ms.methodsToCall["GetByFilter"] = true
	if ms.err != nil {
		return nil, ms.err
	}

	return ms.squad[0:1], nil
}

func (ms *mockSquadStorage) GetOneByID(id string) (*Squad, error) {
	ms.methodsToCall["GetOneByID"] = true
	if ms.err != nil {
		return nil, ms.err
	}
	return ms.squad[0], nil
}

func (ms *mockSquadStorage) UpdateOneByID(id string, updatedSquad Squad) (*Squad, error) {
	ms.methodsToCall["UpdateOneByID"] = true
	if ms.err != nil {
		return nil, ms.err
	}
	return &Squad{
		Id:            ms.squad[0].Id,
		Name:          updatedSquad.Name,
		TeamleadMail:  "john.d@arise.tech",
		Description:   updatedSquad.Description,
		CreatedAt:     ms.squad[0].CreatedAt,
		SkillsRatings: updatedSquad.SkillsRatings,
	}, nil
}

func (ms *mockSquadStorage) DeleteByID(id string) error {
	ms.methodsToCall["DeleteByID"] = true
	return ms.err
}

func (ms *mockSquadStorage) ExpectToCall(methodName string) {
	if ms.methodsToCall == nil {
		ms.methodsToCall = make(map[string]bool)
	}
	ms.methodsToCall[methodName] = false
}

func (ms *mockSquadStorage) Verify(t *testing.T) {
	for methodName, called := range ms.methodsToCall {
		if !called {
			t.Errorf("Expected to call '%s', but it wasn't.", methodName)
		}
	}
}
