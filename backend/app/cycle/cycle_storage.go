package cycle

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"gitdev.devops.krungthai.com/aster/ariskill/app/skill"
	"gitdev.devops.krungthai.com/aster/ariskill/app/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const cycleCollection = "cycles"

type storage struct {
	db *mongo.Database
}

func NewCycleStorage(db *mongo.Database) *storage {
	return &storage{
		db: db,
	}
}

type CycleStorageError struct {
	message string
}

func (e CycleStorageError) Error() string {
	return e.message
}

var invalidRequestError = CycleStorageError{message: "invalid cycle request"}
var cycleNotFoundError = CycleStorageError{message: "cycle not found"}
var dbConnectNotFound = CycleStorageError{message: "cannot connect to mongodb"}

func convertIdToObjectId(id string) (*primitive.ObjectID, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return &objectId, nil
}

// UpdateCycleByID implements Storage.
func (s *storage) UpdateByID(id string, updateCycle Cycle) error {
	cycle := covertCycleDisplayToCycle(CycleDisplay{})
	_ = cycle
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objId}
	update := bson.M{"$set": updateCycle}
	_, err := s.db.Collection(cycleCollection).UpdateOne(ctx, filter, update, options.Update())
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) UpdateByIDSave(id string, updateCycle CycleSave) error {
	cycle := covertCycleDisplayToCycle(CycleDisplay{})
	_ = cycle
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objId}
	update := bson.M{"$set": updateCycle}
	_, err := s.db.Collection(cycleCollection).UpdateOne(ctx, filter, update, options.Update())
	if err != nil {
		return err
	}
	return nil
}

func (s *storage) GetFromUserEmail(email string) ([]*Cycle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"sender_mail": email}
	var cycles []*Cycle
	findOptions := options.Find()

	cursor, err := s.db.Collection(cycleCollection).Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &cycles)
	if err != nil {
		return nil, err
	}

	return cycles, nil
}

// func (s *storage) GetUserAndCycle(email string) ([]*Cycle, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	filter := bson.M{"sender_mail": email}
// 	var cycles []*Cycle
// 	findOptions := options.Find()
// 	//res, err := s.db.Collection(cycleCollection).Find(ctx, filter, findOptions)
// 	cursor, err := s.db.Collection(cycleCollection).Find(ctx, filter, findOptions)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = cursor.All(context.Background(), &cycles)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return cycles, nil
// }

func (s *storage) GetByID(id string) (*Cycle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objId, err := convertIdToObjectId(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objId}

	var cycle Cycle
	err = s.db.Collection(cycleCollection).FindOne(ctx, filter).Decode(&cycle)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, cycleNotFoundError
		}
		return nil, err
	}
	return &cycle, nil
}

// Cycle strut
// ID              primitive.ObjectID `json:"id" bson:"_id"`
// SenderMail      string             `json:"senderMail" bson:"sender_mail"`
// ReceiverMail    string             `json:"receiverMail" bson:"receiver_mail"`
// StartDate       time.Time          `json:"startDate" bson:"start_date"`
// EndDate         time.Time          `json:"endDate" bson:"end_date"`
// QuantitiveSkill []QuantitiveSkill  `json:"quantitiveSkill" bson:"quantitive_skill"`
// IntuitiveSkill  []IntuitiveSkill   `json:"intuitiveSkill" bson:"intuitive_skill"`
// Status          string             `json:"status" bson:"status"`
// Comment         string             `json:"comment" bson:"comment"`

// You have got any problem with this function, contract Mik. 19.51
func (s *storage) InsertOne(cyInput CycleInput, mail string) (*Cycle, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	// Covert CycleInput to Cycle
	cycle := Cycle{SenderMail: mail, ReceiverMail: cyInput.ReceiverMail, StartDate: cyInput.StartDate, EndDate: cyInput.EndDate, QuantitativeSkill: cyInput.QuantitativeSkill, IntuitiveSkill: cyInput.IntuitiveSkill, Status: "Pending", Comment: cyInput.Comment}

	for i, val := range cycle.QuantitativeSkill {
		cycle.QuantitativeSkill[i].LeadGoalScore = val.GoalScore
		fmt.Println(val.GoalScore, cycle.QuantitativeSkill[i].LeadGoalScore)
		cycle.QuantitativeSkill[i].FinalScore = val.PersonalScore
	}
	// Call db to insert data
	res, err := s.db.Collection(cycleCollection).InsertOne(ctx, cycle, options.InsertOne())
	if err != nil {
		return nil, err
	}

	// Get autogenerate ID and return nil if nothing error.
	cycle.ID = res.InsertedID.(primitive.ObjectID)
	return &cycle, nil
}

func (s *storage) GetAllFromReceiverEmail(status string, email string, page int) ([]*Cycle, error) {
	var cycles []*Cycle
	pageSize := 10

	var filter primitive.M
	if status != StatusAll {
		filter = bson.M{"status": status, "receiver_mail": email}
	} else {
		filter = bson.M{"receiver_mail": email}
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64((page - 1) * pageSize))
	findOptions.SetLimit(int64(pageSize))

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := s.db.Collection(cycleCollection).Find(ctx, filter, findOptions)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, cycleNotFoundError
		}
		return nil, dbConnectNotFound
	}

	err = res.All(context.Background(), &cycles)
	if err != nil {
		return nil, err
	}

	// fmt.Println("cycles: ", cycles[0])

	return cycles, nil
}

func (s *storage) DeleteOne(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return invalidRequestError
	}
	filter := bson.M{"_id": oid}

	res, err := s.db.Collection(cycleCollection).DeleteOne(ctx, filter)
	if err != nil {
		return dbConnectNotFound
	}
	if res.DeletedCount == 0 {
		return cycleNotFoundError
	} else {
		return nil
	}
}

const skillCollection = "skills"

func (s *storage) ToDisplayFormat(cy *Cycle) (CycleDisplay, error) {
	ids := []string{}
	for _, val := range cy.QuantitativeSkill {
		ids = append(ids, val.ID.Hex())
	}

	skills, err := getSkillFromCycle(s.db, context.Background(), ids)
	if err != nil {
		return CycleDisplay{}, err
	}

	disPlaySkills := []QuantitativeSkillDisplay{}
	for i, val := range skills {
		displaySkill := QuantitativeSkillDisplay{
			ID:            cy.QuantitativeSkill[i].ID,
			Name:          val.Name,
			Description:   val.Description,
			Logo:          val.Logo,
			PersonalScore: cy.QuantitativeSkill[i].PersonalScore,
			GoalScore:     cy.QuantitativeSkill[i].GoalScore,
			LeadGoalScore: cy.QuantitativeSkill[i].LeadGoalScore,
			FinalScore:    cy.QuantitativeSkill[i].FinalScore,
		}
		disPlaySkills = append(disPlaySkills, displaySkill)
	}
	var cyclesDisplay = CycleDisplay{
		ID:                cy.ID,
		SenderMail:        cy.SenderMail,
		ReceiverMail:      cy.ReceiverMail,
		StartDate:         cy.StartDate,
		EndDate:           cy.EndDate,
		QuantitativeSkill: disPlaySkills,
		IntuitiveSkill:    cy.IntuitiveSkill,
		Status:            cy.Status,
		Comment:           cy.Comment,
	}
	return cyclesDisplay, nil
}

func (s *storage) ToDisplayFormatAll(cycles []*Cycle) []*CycleDisplay {
	allId := [][]string{}

	// skills := [][]skill.Skill{}
	for i := 0; i < len(cycles); i++ {
		ids := []string{}
		for _, val := range cycles[i].QuantitativeSkill {
			ids = append(ids, val.ID.Hex())
		}
		allId = append(allId, ids)
		// skill, _ := sk.GetSkillFromCycle(context.Background(), ids)
		// skills = append(skills, skill)
	}

	skills, err := getAllSkillFromCycle(s.db, context.Background(), allId)
	if err != nil {
		return []*CycleDisplay{}
	}
	cyclesAll := []*CycleDisplay{}

	for i := 0; i < len(skills); i++ {
		toDisplaySkills := []QuantitativeSkillDisplay{}

		for j := 0; j < len(skills[i]); j++ {
			displaySkill := QuantitativeSkillDisplay{
				ID:            cycles[i].QuantitativeSkill[j].ID,
				Name:          skills[i][j].Name,
				Description:   skills[i][j].Description,
				Logo:          skills[i][j].Logo,
				PersonalScore: cycles[i].QuantitativeSkill[j].PersonalScore,
				GoalScore:     cycles[i].QuantitativeSkill[j].GoalScore,
				LeadGoalScore: cycles[i].QuantitativeSkill[j].LeadGoalScore,
				FinalScore:    cycles[i].QuantitativeSkill[j].FinalScore,
			}
			toDisplaySkills = append(toDisplaySkills, displaySkill)
		}
		var cyclesDisplay = CycleDisplay{
			ID:                cycles[i].ID,
			SenderMail:        cycles[i].SenderMail,
			ReceiverMail:      cycles[i].ReceiverMail,
			StartDate:         cycles[i].StartDate,
			EndDate:           cycles[i].EndDate,
			QuantitativeSkill: toDisplaySkills,
			IntuitiveSkill:    cycles[i].IntuitiveSkill,
			Status:            cycles[i].Status,
			State:             cycles[i].State,
			Comment:           cycles[i].Comment,
		}
		cyclesAll = append(cyclesAll, &cyclesDisplay)
	}
	return cyclesAll
}

func getSkillFromCycle(db *mongo.Database, c context.Context, id []string) ([]skill.Skill, error) {
	mapId := map[primitive.ObjectID]int{}
	arrId := []primitive.ObjectID{}
	for _, val := range id {
		oid, _ := primitive.ObjectIDFromHex(val)
		if _, found := mapId[oid]; !found {
			arrId = append(arrId, oid)
		}
	}
	skills := []skill.Skill{}
	filter := bson.M{"_id": bson.M{"$in": arrId}}
	cursor, err := db.Collection(skillCollection).Find(c, filter)

	if err != nil {
		return nil, err
	}

	err = cursor.All(c, &skills)
	if err != nil {
		return nil, err
	}

	allSkill := []skill.Skill{}
	for j := 0; j < len(id); j++ {
		oid, _ := primitive.ObjectIDFromHex(id[j])
		for _, val := range skills {
			if val.ID == oid {
				allSkill = append(allSkill, val)
			}
		}
	}

	return allSkill, nil
}

// GetAllSkillFromCycle implements cycle.skillStorageInterface.
func getAllSkillFromCycle(db *mongo.Database, c context.Context, id [][]string) ([][]skill.Skill, error) {
	// fmt.Println("AAA@@@@@")
	mapId := map[primitive.ObjectID]int{}
	arrId := []primitive.ObjectID{}
	for i := 0; i < len(id); i++ {
		for _, val := range id[i] {
			oid, _ := primitive.ObjectIDFromHex(val)
			if _, found := mapId[oid]; !found {
				arrId = append(arrId, oid)
			}
		}
	}

	skills := []skill.Skill{}
	filter := bson.M{"_id": bson.M{"$in": arrId}}
	cursor, err := db.Collection(skillCollection).Find(c, filter)

	if err != nil {
		return nil, err
	}

	err = cursor.All(c, &skills)
	if err != nil {
		return nil, err
	}

	allSkill := [][]skill.Skill{}
	for i := 0; i < len(id); i++ {
		sk := []skill.Skill{}
		for j := 0; j < len(id[i]); j++ {
			oid, _ := primitive.ObjectIDFromHex(id[i][j])
			for _, val := range skills {
				if val.ID == oid {
					sk = append(sk, val)
				}
			}
		}
		allSkill = append(allSkill, sk)
	}

	return allSkill, nil
}

func covertCycleDisplayToCycle(cycleDisplay CycleDisplay) Cycle {
	quantitiveSK := []QuantitativeSkill{}
	for _, val := range cycleDisplay.QuantitativeSkill {
		skill := QuantitativeSkill{
			ID:            val.ID,
			PersonalScore: val.PersonalScore,
			GoalScore:     val.GoalScore,
			FinalScore:    val.FinalScore,
			Comment:       val.Comment,
		}
		quantitiveSK = append(quantitiveSK, skill)
	}
	cycles := Cycle{
		ID:                cycleDisplay.ID,
		SenderMail:        cycleDisplay.SenderMail,
		ReceiverMail:      cycleDisplay.ReceiverMail,
		StartDate:         cycleDisplay.StartDate,
		EndDate:           cycleDisplay.EndDate,
		QuantitativeSkill: quantitiveSK,
		IntuitiveSkill:    cycleDisplay.IntuitiveSkill,
		Status:            cycleDisplay.Status,
		Comment:           cycleDisplay.Comment,
	}
	return cycles
}

func (s *storage) GetUserDetailWithEmail(cycles []*CycleDisplay) []*CycleWithUserDetail {
	// fmt.Println("cycle:")
	// fmt.Println(cycles[1])
	arrEmail := []string{}
	for _, val := range cycles {
		arrEmail = append(arrEmail, val.SenderMail)
	}
	toSearchEmailArr := []string{}
	toSearchEmail := map[string]int{}
	for _, val := range arrEmail {
		if _, found := toSearchEmail[val]; !found {
			toSearchEmailArr = append(toSearchEmailArr, val)
			toSearchEmail[val] = 1
		}
	}

	users := []user.User{}
	filter := bson.M{"email": bson.M{"$in": toSearchEmailArr}}
	cursor, err := s.db.Collection("users").Find(context.Background(), filter)

	if err != nil {
		panic(err)
	}

	err = cursor.All(context.Background(), &users)
	if err != nil {
		panic(err)
	}
	// fmt.Println(toSearchEmailArr, users)

	cyclesAll := []*CycleWithUserDetail{}

	for i := 0; i < len(arrEmail); i++ {
		for _, val := range users {
			if val.Email == arrEmail[i] {
				cycle := &CycleWithUserDetail{
					ID:                cycles[i].ID,
					UID:               val.ID,
					FirstName:         val.FirstName,
					LastName:          val.LastName,
					JobRole:           val.JobRole,
					Level:             val.Level,
					SenderMail:        cycles[i].SenderMail,
					ReceiverMail:      cycles[i].ReceiverMail,
					StartDate:         cycles[i].StartDate,
					EndDate:           cycles[i].EndDate,
					QuantitativeSkill: cycles[i].QuantitativeSkill,
					IntuitiveSkill:    cycles[i].IntuitiveSkill,
					Status:            cycles[i].Status,
					State:             cycles[i].State,
					Comment:           cycles[i].Comment,
				}
				cyclesAll = append(cyclesAll, cycle)
				break
			}
		}
	}

	// fmt.Println("cycles all: ", cyclesAll[0])
	return cyclesAll
}

func (s *storage) UpdateUserFinalScore(id string) error {
	cycle, err := s.GetByID(id)
	if err != nil {
		return err
	}

	response := GetFinalScoreResponse{
		Email: cycle.SenderMail,
	}
	for _, skill := range cycle.QuantitativeSkill {
		response.QuantiativeSkills = append(response.QuantiativeSkills,
			QuantiativeSkillToUpdate{
				ID:         skill.ID,
				FinalScore: skill.FinalScore,
			})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	filter := bson.M{"email": response.Email}
	update := bson.M{}

	var arrayFilters []interface{}
	for i, skill := range response.QuantiativeSkills {
		filterIdentifier := fmt.Sprintf("elem%d", i)
		arrayFilters = append(arrayFilters, bson.M{filterIdentifier + ".skillID": skill.ID})
	}
	if len(arrayFilters) > 0 {
		update["$set"] = bson.M{}
		for i, skill := range response.QuantiativeSkills {
			filterIdentifier := fmt.Sprintf("elem%d", i)
			update["$set"].(bson.M)["technical_skills.$["+filterIdentifier+"].score"] = skill.FinalScore
		}
		updateOptions := options.Update().SetArrayFilters(options.ArrayFilters{Filters: arrayFilters})
		_, err = s.db.Collection("users").UpdateOne(ctx, filter, update, updateOptions)
		if err != nil {
			return err
		}
	}
	for _, skill := range response.QuantiativeSkills {
		addFilter := bson.M{"email": response.Email, "technical_skills.skillID": bson.M{"$ne": skill.ID}}
		addToSet := bson.M{"$addToSet": bson.M{"technical_skills": bson.M{"skillID": skill.ID, "score": skill.FinalScore}}}
		_, err = s.db.Collection("users").UpdateOne(ctx, addFilter, addToSet)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateCycleStatus(status string) bool {
	switch status {
	case StatusAll, StatusPending, StatusApproved, StatusRunning, StatusDone:
		return true
	}

	return false
}

// To improve
func (s *storage) ToDisplayFormatAllExpermentConCurrency(cycles []*Cycle) []*CycleDisplay {
	var wg sync.WaitGroup
	c := make(chan *CycleDisplay, len(cycles))
	temp := func(cy *Cycle, s *storage, c chan *CycleDisplay, wg *sync.WaitGroup) {
		t, err := s.ToDisplayFormat(cy)
		if err != nil {
			//TODO: what should we do
			slog.Error(err.Error())
		}
		c <- &t
		wg.Done()
	}

	for i := 0; i < len(cycles); i++ {
		wg.Add(1)
		go temp(cycles[i], s, c, &wg)
	}
	wg.Wait()
	cyclesAll := []*CycleDisplay{}
	for v := range c {
		cyclesAll = append(cyclesAll, v)
	}
	return cyclesAll
}

// New CyCle Storage

const newCycleCollection = "new_cycles"
const userCollection = "users"

func (s *storage) GetNewByID(id string) (*NewCycle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objId, err := convertIdToObjectId(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objId}

	var cycle NewCycle
	err = s.db.Collection(newCycleCollection).FindOne(ctx, filter).Decode(&cycle)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, cycleNotFoundError
		}
		return nil, err
	}
	return &cycle, nil
}

func toHardSkillDisplay(hs *[]HardSkill) []HardSkillDisplay {
	hsd := []HardSkillDisplay{}
	for _, val := range *hs {
		newHsd := &HardSkillDisplay{
			ID:            val.ID,
			Name:          val.Name,
			Description:   val.Description,
			SkillLevels:   val.SkillLevels,
			PersonalScore: val.PersonalScore,
			GoalScore:     val.GoalScore,
		}

		hsd = append(hsd, *newHsd)
	}

	return hsd
}

func (s *storage) ToNewUserDetailFormat(cy *NewCycle) (*NewCycleWithUserDetail, error) {
	filter := bson.M{"email": cy.AriserMail}
	result := s.db.Collection("users").FindOne(context.Background(), filter)

	var user user.User
	err := result.Decode(&user)
	if err != nil {
		return nil, err
	}

	cyd := &NewCycleWithUserDetail{
		ID:             cy.ID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		JobRole:        user.JobRole,
		Level:          user.Level,
		TeamLeaderMail: cy.TeamLeaderMail,
		AriserMail:     cy.AriserMail,
		StartDate:      cy.StartDate,
		EndDate:        cy.EndDate,
		Status:         cy.Status,
		Comment:        cy.Comment,
		HardSkills:     toHardSkillDisplay(&cy.HardSkills),
		State:          cy.State,
	}

	return cyd, nil
}

func (s *storage) GetLatestCycleFromUserEmail(email string) (*NewCycle, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"ariserMail": email, "status": "In Progress"}
	var cycles *NewCycle
	findOptions := options.FindOne().SetSort(bson.M{"startDate": -1})
	err := s.db.Collection(newCycleCollection).FindOne(ctx, filter, findOptions).Decode(&cycles)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return cycles, nil
}

func (s *storage) GetUsersHardSkillByEmail(ctx context.Context, email string) (*user.User, error) {
	var userData user.User
	filter := bson.M{"email": email}
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := s.db.Collection(userCollection).FindOne(ctx, filter).Decode(&userData); err != nil {
		return nil, err
	}
	return &userData, nil
}

func (s *storage) UpdateHardSkillsByEmail(ctx context.Context, email string, goalSkillRequest UpdateGoalSkillsRequest) (*NewCycle, error) {
	userDetail, err := s.GetUsersHardSkillByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	cycles, err := s.GetLatestCycleFromUserEmail(email)
	if err != nil {
		return nil, err
	}
	// use mapping to map data from userDetail
	userMapping := make(map[string]user.MyHardSkill)
	for _, v := range userDetail.HardSkills {
		userMapping[v.Name] = v
	}

	// loop data from mapping to hardskill Request
	for i, value := range goalSkillRequest.HardSkills {
		goalSkillRequest.HardSkills[i].PersonalScore = userMapping[value.Name].CurrentLevel
		if goalSkillRequest.HardSkills[i].PersonalScore+1 < goalSkillRequest.HardSkills[i].GoalScore || goalSkillRequest.HardSkills[i].PersonalScore > goalSkillRequest.HardSkills[i].GoalScore {
			return nil, fmt.Errorf("miss match goal-score")
		}
	}

	cycles.HardSkills = goalSkillRequest.HardSkills
	cycles.Status = "Pending"
	cycles.State = "Review"

	filter := bson.M{"_id": primitive.ObjectID(cycles.ID)}

	update := bson.M{"$set": cycles}

	_, err = s.db.Collection(newCycleCollection).UpdateOne(ctx, filter, update, options.Update())

	if err != nil {
		return nil, err
	}

	return cycles, nil
}
