package squad

import (
	"context"
	"io"
	"reflect"

	"gitdev.devops.krungthai.com/aster/ariskill/app"
	"gitdev.devops.krungthai.com/aster/ariskill/app/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type squadStorageInterface interface {
	InsertOne(userId string, squadToInsert Squad) (*Squad, error)
	GetAll() ([]*Squad, error)
	GetByFilter(filter regexFilter) ([]*Squad, error)
	GetOneByID(id string) (*Squad, error)
	UpdateOneByID(id string, updatedSquad Squad) (*Squad, error)
	DeleteByID(id string) error
}

type memberStorageInterface interface {
	GetAllBySquadId(context context.Context, squadId primitive.ObjectID) ([]user.User, error)
}

type regexFilter map[string]map[string]primitive.Regex

type squadHandler struct {
	storage       squadStorageInterface
	memberStorage memberStorageInterface
}

func NewSquadHandler(storage squadStorageInterface, memberStorage memberStorageInterface) *squadHandler {
	return &squadHandler{
		storage:       storage,
		memberStorage: memberStorage,
	}
}

type squadHandlerError struct {
	message string
}

func (e squadHandlerError) Error() string {
	return e.message
}

var invalidFilterError = squadHandlerError{message: "invalid filter"}
var invalidSquadInputError = squadHandlerError{message: "invalid squad input"}
var missingSkillsRatingsError = squadHandlerError{message: "missing skills ratings, must be at least 3"}
var missingRatingsError = squadHandlerError{message: "missing ratings, must be at least 1"}

var squadMemberNotFoundError = squadHandlerError{message: "squad member not found"}

// UpdateOneByID godoc
//
//	@summary		UpdateOneByID
//	@description	Update a squad when lead edit squad data
//	@tags			squad
//	@id				UpdateOneByID
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@param			id		path		string			true	"Squad ID"
//	@param			squad	body		squad.Squad		true	"Squad"
//	@response		200		{object}	squad.Squad		"OK"
//	@response		400		{object}	app.Response	"Bad Request"
//	@response		401		{object}	app.Response	"Unauthorized"
//	@response		404		{object}	app.Response	"Not Found"
//	@response		500		{object}	app.Response	"Internal Server Error"
//	@router			/squads/{squadID} [put]
func (h *squadHandler) UpdateOneByID(c app.Context) {
	sqId := c.Param("squadID")
	oldSquad, err := h.storage.GetOneByID(sqId)
	if err != nil {
		if err == invalidIdError {
			c.BadRequest(err)
			return
		}
		if err == squadNotFoundError {
			c.NotFound(err)
			return
		}

		c.InternalServerError(err)
		return
	}

	var updatedSquad Squad
	err = c.Bind(&updatedSquad)
	if err != nil {
		c.BadRequest(invalidSquadInputError)
		return
	}

	combinedSquad := *oldSquad
	updatedSquadValue := reflect.ValueOf(updatedSquad)
	updatedSquadType := updatedSquadValue.Type()
	for i := 0; i < updatedSquadValue.NumField(); i++ {
		fieldValue := updatedSquadValue.Field(i)
		fieldName := updatedSquadType.Field(i).Name

		combinedSquadValue := reflect.ValueOf(&combinedSquad).Elem()
		if !combinedSquadValue.FieldByName(fieldName).IsValid() {
			c.BadRequest(invalidSquadInputError)
			return
		}

		if !fieldValue.IsZero() {
			combinedSquadValue.FieldByName(fieldName).Set(fieldValue)
		}
	}

	res, err := h.storage.UpdateOneByID(sqId, combinedSquad)
	if err != nil {
		if err == invalidIdError {
			c.BadRequest(err)
		} else {
			c.InternalServerError(err)
		}
		return
	}

	c.OK(res)
}

// GetAll godoc
//
//	@summary		GetAll
//	@description	Get all squads
//	@tags			squad
//	@id				GetAll
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@param			squadFilter	body		SquadFilter		true	"Filter for squad"
//	@response		200			{array}		squad.Squad		"OK"
//	@response		400			{object}	app.Response	"Bad Request"
//	@response		401			{object}	app.Response	"Unauthorized"
//	@response		404			{object}	app.Response	"Not Found"
//	@response		500			{object}	app.Response	"Internal Server Error"
//	@router			/squads [get]
func (h *squadHandler) GetAll(c app.Context) {
	filter := SquadFilter{}
	err := c.Bind(&filter)
	if err != nil && err != io.EOF {
		c.BadRequest(invalidFilterError)
		return
	}
	var res []*Squad
	if err == nil {
		regexFilter := regexFilter{}
		if filter.Name != "" {
			regexFilter["name"] = map[string]primitive.Regex{"$regex": {Pattern: filter.Name + ".*", Options: "i"}}
		}
		if filter.Description != "" {
			regexFilter["desc"] = map[string]primitive.Regex{"$regex": {Pattern: filter.Description + ".*", Options: "i"}}
		}
		res, err = h.storage.GetByFilter(regexFilter)
	} else {
		res, err = h.storage.GetAll()
	}

	if err != nil {
		if err == invalidIdError {
			c.BadRequest(err)
			return
		}
		if err == squadNotFoundError {
			c.NotFound(err)
			return
		}

		c.InternalServerError(err)
		return
	}

	c.OK(res)
}

// GetOneByID godoc
//
//	@summary		GetOneByID
//	@description	Get one squad by squad id
//	@tags			squad
//	@id				GetOneByID
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@param			id	path		string			true	"Squad ID"
//	@response		200	{object}	squad.Squad		"OK"
//	@response		400	{object}	app.Response	"Bad Request"
//	@response		401	{object}	app.Response	"Unauthorized"
//	@response		500	{object}	app.Response	"Internal Server Error"
//	@router			/squads/{id} [get]
func (h *squadHandler) GetOneByID(c app.Context) {
	id := c.Param("squadID")
	res, err := h.storage.GetOneByID(id)
	if err != nil {
		if err == invalidIdError {
			c.BadRequest(err)
			return
		}

		if err == squadNotFoundError {
			c.NotFound(err)
			return
		}

		c.InternalServerError(err)
		return
	}

	c.OK(res)
}

// InsertOneByID godoc
//
//	@summary		InsertOneByID
//	@description	Create a squad
//	@tags			squad
//	@id				InsertOneByID
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@param			insert	body		squad.Squad		true	"Squad object"
//	@response		200		{object}	squad.Squad		"OK"
//	@response		400		{object}	app.Response	"Bad Request"
//	@response		401		{object}	app.Response	"Unauthorized"
//	@response		500		{object}	app.Response	"Internal Server Error"
//	@router			/squads [post]
func (handler *squadHandler) InsertOneByID(c app.Context) {
	uid := c.GetString("profileID")
	var insert Squad
	err := c.Bind(&insert)
	if err != nil {
		c.BadRequest(invalidSquadInputError)
		return
	}

	// Check that there's at least 3 skills in the squad
	if len(insert.SkillsRatings) < 3 {
		c.BadRequest(missingSkillsRatingsError)
		return
	}

	// Each skill already has 1 rating
	for _, v := range insert.SkillsRatings {
		if len(v.Ratings) < 1 {
			c.BadRequest(missingRatingsError)
			return
		}
	}

	res, err := handler.storage.InsertOne(uid, insert)
	if err != nil {
		c.InternalServerError(err)
		return
	}

	// TODO: after succesfully inserted, added creator as role owner
	c.OK(res)
}

// 	c.OK(res)
// }

// DeleteSquadByID godoc
//
//	@summary		DeleteSquadByID
//	@description	Delete a squad
//	@tags			squad
//	@id				DeleteSquadByID
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@param			id	path		string			true	"Squad ID"
//	@response		200	{object}	squad.Squad		"OK"
//	@response		400	{object}	app.Response	"Bad Request"
//	@response		401	{object}	app.Response	"Unauthorized"
//	@response		500	{object}	app.Response	"Internal Server Error"
//	@router			/squads/{id} [delete]
func (handler *squadHandler) DeleteByID(c app.Context) {
	id := c.Param("squadID")

	err := handler.storage.DeleteByID(id)
	if err != nil {
		if err == invalidIdError {
			c.BadRequest(err)
			return
		}

		if err == squadNotFoundError {
			c.NotFound(err)
			return
		}

		c.InternalServerError(err)
		return
	}

	c.OK(nil)
}

// CalculateSquadMemberAveragePerSkill godoc
//
//	@summary		CalculateSquadMemberAveragePerSkill
//	@description	Calculate squad member average per skill
//	@tags			squad
//	@id				CalculateSquadMemberAveragePerSkill
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@param			squadID	path		string					true	"Squad ID"
//	@response		200		{object}	SquadAverageSkillOutput	"OK"
//	@response		400		{object}	app.Response			"Bad Request"
//	@response		401		{object}	app.Response			"Unauthorized"
//	@response		404		{object}	app.Response			"Not found"
//	@response		500		{object}	app.Response			"Internal Server Error"
//	@router			/squads/{squadID}/member-skills-avg [get]
func (handler *squadHandler) CalculateSquadMemberAveragePerSkill(c app.Context) {
	id := c.Param("squadID")
	if id == "" {
		c.BadRequest(invalidIdError)
		return
	}

	res, err := handler.storage.GetOneByID(id)
	if err != nil {
		if err == invalidIdError {
			c.BadRequest(err)
			return
		}

		if err == squadNotFoundError {
			c.NotFound(err)
			return
		}

		c.InternalServerError(err)
		return
	}
	users, err := handler.memberStorage.GetAllBySquadId(context.Background(), res.Id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.NotFound(squadMemberNotFoundError)
			return
		}

		c.InternalServerError(err)
		return
	}

	averagesPerSkill := SquadAverageSkillOutput{
		SquadId:       res.Id,
		AveragesSkill: []AverageSkill{},
	}
	for _, skill := range res.SkillsRatings {
		sum := 0
		total := 0
		for _, user := range users {
			userSkills := user.TechnicalSkill
			for _, userSkill := range userSkills {
				if skill.SkillId == userSkill.SkillID {
					sum += userSkill.Score
					total += 1
				}
			}
		}
		average := float64(sum) / float64(total)

		newAverageSkill := AverageSkill{
			SkillId: skill.SkillId,
			Average: average,
		}
		averagesPerSkill.AveragesSkill = append(averagesPerSkill.AveragesSkill, newAverageSkill)
	}
	c.OK(averagesPerSkill)
}

// GetAvgSkillRatingByID godoc
//
//	@summary		GetAvgSkillRatingByID
//	@description	Get one squad with average skill ratings
//	@tags			squad
//	@id				GetAvgSkillRatingByID
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@param			squadID	path		string					true	"Squad ID"
//	@response		200		{object}	SquadAverageSkillOutput	"OK"
//	@response		400		{object}	app.Response			"Bad Request"
//	@response		401		{object}	app.Response			"Unauthorized"
//	@response		405		{object}	app.Response			"Store Error"
//	@response		500		{object}	app.Response			"Internal Server Error"
//	@router			/squads/{squadID}/skills-require-avg [get]
func (handler *squadHandler) GetAvgSkillRatingByID(c app.Context) {
	sqid := c.Param("squadID") // Get from url param
	sq, err := handler.storage.GetOneByID(sqid)
	if err != nil {
		if err == invalidIdError {
			c.BadRequest(err)
			return
		}
		if err == squadNotFoundError {
			c.NotFound(err)
			return
		}

		c.InternalServerError(err)
		return
	}

	res := &SquadAverageSkillOutput{
		SquadId: sq.Id,
	}

	// Loop over each skill
	// Find average
	for i := 0; i < len(sq.SkillsRatings); i++ {
		avgScore := 0
		if len(sq.SkillsRatings[i].Ratings) > 0 {
			total := 0
			for j := 0; j < len(sq.SkillsRatings[i].Ratings); j++ {
				total += sq.SkillsRatings[i].Ratings[j].Score
			}
			avgScore = total * 10 / len(sq.SkillsRatings[i].Ratings)
		}
		res.AveragesSkill = append(res.AveragesSkill, AverageSkill{
			SkillId: sq.SkillsRatings[i].SkillId,
			Average: float64(avgScore) / 10,
		})
	}

	// Remove the embedded skills to reduce memory usage and aid in testing
	// res.Squad.Skills = nil

	c.OK(res)
}
