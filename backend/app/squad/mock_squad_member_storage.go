package squad

import (
	"context"
	"testing"

	"gitdev.devops.krungthai.com/aster/ariskill/app/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockUserStorage struct {
	users         []user.User
	methodsToCall map[string]bool
	err           error
}

func (m *mockUserStorage) GetAllBySquadId(context context.Context, squadId primitive.ObjectID) ([]user.User, error) {
	m.methodsToCall["GetAllBySquadId"] = true
	if m.err != nil {
		return nil, m.err
	}

	return m.users, nil
}

func (ms *mockUserStorage) ExpectToCall(methodName string) {
	if ms.methodsToCall == nil {
		ms.methodsToCall = make(map[string]bool)
	}
	ms.methodsToCall[methodName] = false
}

func (ms *mockUserStorage) Verify(t *testing.T) {
	for methodName, called := range ms.methodsToCall {
		if !called {
			t.Errorf("Expected to call '%s', but it wasn't.", methodName)
		}
	}
}
