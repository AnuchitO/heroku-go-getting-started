package squad

import (
	"context"

	"gitdev.devops.krungthai.com/aster/ariskill/app/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *mockSquadStorage) GetAllBySquadId(context context.Context, squadId primitive.ObjectID) ([]user.User, error) {
	m.methodsToCall["GetAllBySquadId"] = true
	if m.err != nil {
		return nil, m.err
	}

	return m.users, nil
}

// func (ms *mockSquadStorage) ExpectToCall(methodName string) {
// 	if ms.methodsToCall == nil {
// 		ms.methodsToCall = make(map[string]bool)
// 	}
// 	ms.methodsToCall[methodName] = false
// }

// func (ms *mockSquadStorage) Verify(t *testing.T) {
// 	for methodName, called := range ms.methodsToCall {
// 		if !called {
// 			t.Errorf("Expected to call '%s', but it wasn't.", methodName)
// 		}
// 	}
// }
