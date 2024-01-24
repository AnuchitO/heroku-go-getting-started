package cycle

import (
	"context"
	"fmt"
	"testing"

	"gitdev.devops.krungthai.com/aster/ariskill/app/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockCycleStorage struct {
	cycles       []*Cycle
	cycle        *Cycle
	cyclesReturn []*Cycle
	// newCycle     *NewCycle
	newCycles []*NewCycle

	methodsToCall map[string]bool
	err           error
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

// GetAllFromEmail implements Storage.
func (ms *mockCycleStorage) GetAllFromReceiverEmail(status string, email string, page int) ([]*Cycle, error) {
	ms.methodsToCall["GetAllFromEmail"] = true
	if ms.err != nil {
		return nil, ms.err
	}

	return ms.cyclesReturn, nil
}

// GetAll implements Storage.
func (m *mockCycleStorage) GetAll() ([]*Cycle, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.cycles, nil
}

func (m *mockCycleStorage) DeleteOne(id string) error {
	m.methodsToCall["DeleteByID"] = true
	return m.err
}

// InsertOne implements Storage.
func (m *mockCycleStorage) InsertOne(cy CycleInput, mail string) (*Cycle, error) {
	if m.err != nil {
		return nil, m.err
	}
	m.cycle.IntuitiveSkill = []IntuitiveSkill{}
	return m.cycle, nil
}

// UpdateCycleByID implements Storage.
func (m *mockCycleStorage) UpdateByID(id string, updateCycle Cycle) error {
	m.methodsToCall["UpdateByID"] = true
	return m.err
}

func (m *mockCycleStorage) UpdateByIDSave(id string, updateCycle CycleSave) error {
	m.methodsToCall["UpdateByIDSave"] = true
	return m.err
}

func (ms *mockCycleStorage) GetByID(id string) (*Cycle, error) {
	ms.methodsToCall["GetByID"] = true
	if ms.err != nil {
		return nil, ms.err
	}
	return ms.cycles[0], nil
}

func (ms *mockCycleStorage) ExpectToCall(methodName string) {
	if ms.methodsToCall == nil {
		ms.methodsToCall = make(map[string]bool)
	}
	ms.methodsToCall[methodName] = false
}

func (ms *mockCycleStorage) Verify(t *testing.T) {
	for methodName, called := range ms.methodsToCall {
		if !called {
			t.Errorf("Expected to call '%s', but it wasn't.", methodName)
		}
	}
}

func (ms *mockCycleStorage) ToDisplayFormat(cy *Cycle) (CycleDisplay, error) {
	cycleDisplayMock := CycleDisplay{
		ID:                cy.ID,
		SenderMail:        cy.SenderMail,
		ReceiverMail:      cy.ReceiverMail,
		StartDate:         cy.StartDate,
		EndDate:           cy.EndDate,
		QuantitativeSkill: []QuantitativeSkillDisplay{},
		IntuitiveSkill:    cy.IntuitiveSkill,
		Status:            cy.Status,
		Comment:           cy.Comment,
	}
	return cycleDisplayMock, nil
}

func (ms *mockCycleStorage) ToDisplayFormatAll(cycles []*Cycle) []*CycleDisplay {
	cyclesDisplayMock := []*CycleDisplay{}
	for _, val := range cycles {
		cycleDisplayMock := CycleDisplay{
			ID:                val.ID,
			SenderMail:        val.SenderMail,
			ReceiverMail:      val.ReceiverMail,
			StartDate:         val.StartDate,
			EndDate:           val.EndDate,
			QuantitativeSkill: []QuantitativeSkillDisplay{},
			IntuitiveSkill:    val.IntuitiveSkill,
			Status:            val.Status,
			Comment:           val.Comment,
		}
		cyclesDisplayMock = append(cyclesDisplayMock, &cycleDisplayMock)
	}
	return cyclesDisplayMock
}

func (ms *mockCycleStorage) ToDisplayFormatAllExpermentConCurrency(cycles []*Cycle) []*CycleDisplay {
	return nil
}

func (ms *mockCycleStorage) GetFromUserEmail(email string) ([]*Cycle, error) {
	ms.methodsToCall["GetFromUserEmail"] = true
	if ms.err != nil {
		return nil, ms.err
	}

	return ms.cycles, nil
}

func (ms *mockCycleStorage) GetUserDetailWithEmail(cycles []*CycleDisplay) []*CycleWithUserDetail {
	if ms.cyclesReturn != nil {
		cyclesAll := []*CycleWithUserDetail{}
		for i := 0; i < len(cycles); i++ {
			cycle := &CycleWithUserDetail{
				ID:                cycles[i].ID,
				UID:               primitive.NewObjectID().Hex(),
				FirstName:         "Ariser",
				LastName:          "by krungthai",
				JobRole:           "full-stack",
				Level:             "junior",
				SenderMail:        cycles[i].SenderMail,
				ReceiverMail:      cycles[i].ReceiverMail,
				StartDate:         cycles[i].StartDate,
				EndDate:           cycles[i].EndDate,
				QuantitativeSkill: cycles[i].QuantitativeSkill,
				IntuitiveSkill:    cycles[i].IntuitiveSkill,
				Status:            cycles[i].Status,
				Comment:           cycles[i].Comment,
			}
			cyclesAll = append(cyclesAll, cycle)
		}

		return cyclesAll
	}
	return nil
}

func (ms *mockCycleStorage) UpdateUserFinalScore(id string) error {
	ms.methodsToCall["UpdateUserFinalScore"] = true
	return ms.err
}

func (ms *mockCycleStorage) GetNewByID(id string) (*NewCycle, error) {
	ms.methodsToCall["GetNewByID"] = true
	if ms.err != nil {
		return nil, ms.err
	}
	return ms.newCycles[0], nil
}

func (ms *mockCycleStorage) ToNewUserDetailFormat(cy *NewCycle) (*NewCycleWithUserDetail, error) {
	cycle := &NewCycleWithUserDetail{
		ID:             cy.ID,
		FirstName:      "Ariser",
		LastName:       "by krungthai",
		JobRole:        "full-stack",
		Level:          "junior",
		AriserMail:     cy.AriserMail,
		TeamLeaderMail: cy.TeamLeaderMail,
		StartDate:      cy.StartDate,
		EndDate:        cy.EndDate,
		HardSkills:     []HardSkillDisplay{},
		Status:         cy.Status,
		Comment:        cy.Comment,
	}
	return cycle, nil
}

func (ms *mockCycleStorage) UpdateHardSkillsByEmail(ctx context.Context, email string, goalSkillRequest UpdateGoalSkillsRequest) (*NewCycle, error) {
	panic("not Implement")
}

func (ms *mockCycleStorage) GetLatestCycleFromUserEmail(email string) (*NewCycle, error) {
	panic("not Implement")
}

// TODO : Uncomment to Test add hardSkills a Cycle to project
// new cycle
type mockNewCycleStorage struct {
	newCycle *NewCycle
	user     user.User
	err      error
}

func (ms *mockNewCycleStorage) UpdateByID(id string, updateCycle Cycle) error {
	panic("not Implement")
}
func (ms *mockNewCycleStorage) UpdateByIDSave(id string, updateCycle CycleSave) error {
	panic("not Implement")
}
func (ms *mockNewCycleStorage) InsertOne(cy CycleInput, mail string) (*Cycle, error) {
	panic("not Implement")
}
func (ms *mockNewCycleStorage) GetAllFromReceiverEmail(status string, email string, page int) ([]*Cycle, error) {
	panic("not Implement")
}
func (ms *mockNewCycleStorage) UpdateUserFinalScore(id string) error {
	panic("not Implement")
}
func (ms *mockNewCycleStorage) GetFromUserEmail(email string) ([]*Cycle, error) {
	panic("not Implement")
}
func (ms *mockNewCycleStorage) GetByID(id string) (*Cycle, error) {
	panic("not Implement")
}
func (ms *mockNewCycleStorage) DeleteOne(id string) error {
	panic("not Implement")
}
func (ms *mockNewCycleStorage) ToDisplayFormat(cy *Cycle) (CycleDisplay, error) {
	panic("not Implement")
}
func (ms *mockNewCycleStorage) ToDisplayFormatAll(cycles []*Cycle) []*CycleDisplay {
	panic("not Implement")
}
func (ms *mockNewCycleStorage) GetUserDetailWithEmail(cycles []*CycleDisplay) []*CycleWithUserDetail {
	panic("not Implement")
}
func (ms *mockNewCycleStorage) UpdateHardSkillsByEmail(ctx context.Context, email string, goalSkillRequest UpdateGoalSkillsRequest) (*NewCycle, error) {
	if ms.err != nil {
		return nil, ms.err
	}
	ms.newCycle.HardSkills = goalSkillRequest.HardSkills
	return ms.newCycle, nil
}

func (ms *mockNewCycleStorage) GetLatestCycleFromUserEmail(email string) (*NewCycle, error) {
	if ms.err != nil {
		return nil, ms.err
	}
	return ms.newCycle, nil
}

func (ms *mockNewCycleStorage) GetUsersHardSkillByEmail(ctx context.Context, email string) (*user.User, error) {
	if ms.err != nil {
		return nil, ms.err
	}
	return &ms.user, nil
}

func (ms *mockNewCycleStorage) GetNewByID(id string) (*NewCycle, error) {
	panic("not Implement")
}

func (ms *mockNewCycleStorage) ToNewUserDetailFormat(cy *NewCycle) (*NewCycleWithUserDetail, error) {
	panic("not Implement")
}
