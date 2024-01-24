package cycle

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"gitdev.devops.krungthai.com/aster/ariskill/app"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Storage interface {
	UpdateByID(id string, updateCycle Cycle) error
	UpdateByIDSave(id string, updateCycle CycleSave) error
	InsertOne(cy CycleInput, mail string) (*Cycle, error)
	GetAllFromReceiverEmail(status string, email string, page int) ([]*Cycle, error)
	UpdateUserFinalScore(id string) error
	GetFromUserEmail(email string) ([]*Cycle, error)
	GetByID(id string) (*Cycle, error)
	DeleteOne(id string) error
	ToDisplayFormat(cy *Cycle) (CycleDisplay, error)
	ToDisplayFormatAll(cycles []*Cycle) []*CycleDisplay
	GetUserDetailWithEmail(cycles []*CycleDisplay) []*CycleWithUserDetail
	GetNewByID(id string) (*NewCycle, error)
	ToNewUserDetailFormat(cy *NewCycle) (*NewCycleWithUserDetail, error)
	GetLatestCycleFromUserEmail(email string) (*NewCycle, error)
	UpdateHardSkillsByEmail(ctx context.Context, email string, goalSkillRequest UpdateGoalSkillsRequest) (*NewCycle, error)
	// UpdateByEmail(email string, updateCycle UpdateGoalSkillsRequest) error
	// GetFromUserEmailNewCycle(email string) (*NewCycle, error)
}

type cycleHandler struct {
	storage Storage
}

// type userStorageInterface interface {
// 	GetUserById(ctx context.Context, id string) (user.User, error)
// }

type cycleHandlerError struct {
	message string
}

// Error implements error.
func (e cycleHandlerError) Error() string {
	return e.message
}

var numberOfSkillError = cycleHandlerError{message: "each type of skill should have atleast one"}
var invalidInsertOneInputError = cycleHandlerError{message: "invalid or missing required field"}

func NewCycleHandler(st Storage) *cycleHandler {
	return &cycleHandler{
		storage: st,
	}
}

// InsertOne godoc
//
//	@summary		InsertOne
//	@description	Insert new Cycles
//	@tags			cycle
//	@id				InsertOneCycle
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@param			reqJson	body		CycleInput			true	"Cycle input Object"
//	@response		200		{array}		cycle.CycleDisplay	"OK"
//	@response		400		{object}	app.Response		"Bad Request"
//	@response		401		{object}	app.Response		"Unauthorized"
//	@response		404		{object}	app.Response		"Not Found"
//	@response		500		{object}	app.Response		"Internal Server Error"
//	@router			/cycles [post]
func (h *cycleHandler) InsertOne(c app.Context) {
	mail := c.GetString("email")
	var insert CycleInput

	// It'll return EOF error if your POST with wrong BODY format.
	err := c.Bind(&insert)
	if err != nil {
		c.BadRequest(invalidInsertOneInputError)
		return
	}

	// User have to select atleast 1 quantitive skill.
	if len(insert.QuantitativeSkill) < 1 && len(insert.IntuitiveSkill) < 1 {
		c.BadRequest(numberOfSkillError)
		return
	}

	res, err := h.storage.InsertOne(insert, mail)

	if err != nil {
		c.InternalServerError(err)
		return
	}

	toDisplay, err := h.storage.ToDisplayFormat(res)
	if err != nil {
		c.InternalServerError(err)
		return
	}

	c.OK(toDisplay)
}

// UpdateCycle godoc
//
//	@summary		UpdateByID
//	@description	Update cycle with specific ID
//	@tags			cycle
//	@id				UpdateByID
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@param			reqJson	body		Cycle				true	"Cycle input Object"
//	@response		200		{array}		cycle.CycleDisplay	"OK"
//	@response		400		{object}	app.Response		"Bad Request"
//	@response		401		{object}	app.Response		"Unauthorized"
//	@response		404		{object}	app.Response		"Not Found"
//	@response		500		{object}	app.Response		"Internal Server Error"
//	@router			/cycles/{id}  [post]
func (cy *cycleHandler) UpdateByID(c app.Context) {
	id := c.Param("id")
	var json Cycle

	err := c.ShouldBindJSON(&json)

	if err != nil {
		c.BadRequest(err)
		return
	}
	err = cy.storage.UpdateByID(id, json)

	if err != nil {
		c.InternalServerError(err)
		return
	}
	// TODO: Implement function GetByID
	res, err := cy.storage.GetByID(id)
	if err != nil {
		c.InternalServerError(err)
		return
	}

	toDisplay, err := cy.storage.ToDisplayFormat(res)
	if err != nil {
		c.InternalServerError(err)
		return
	}

	c.OK(toDisplay)
}

func (cy *cycleHandler) UpdateByIDSave(c app.Context) {
	id := c.Param("id")
	var json CycleSave
	// var json map[string]interface{}

	err := c.ShouldBindJSON(&json)

	if err != nil {
		c.BadRequest(err)
		return
	}
	// delete(json, "status")

	err = cy.storage.UpdateByIDSave(id, json)

	if err != nil {
		c.InternalServerError(err)
		return
	}
	// TODO: Implement function GetByID
	res, _ := cy.storage.GetByID(id)
	toDisplay, _ := cy.storage.ToDisplayFormat(res)
	c.OK(toDisplay)
}

// GetAllFromReceiverEmail godoc
//
//	@summary		GetAllFromReceiverEmail
//	@description	Get Cycles by sending status and page and receiver email (in context)
//	@tags			cycle
//	@id				GetAllFromReceiverEmail
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@response		200	{array}		cycle.CycleDisplay	"OK"
//	@response		400	{object}	app.Response		"Bad Request"
//	@response		401	{object}	app.Response		"Unauthorized"
//	@response		404	{object}	app.Response		"Not Found"
//	@response		500	{object}	app.Response		"Internal Server Error"
//	@router			/cycles/email/{status}/{page} [get]
func (cy *cycleHandler) GetAllFromReceiverEmail(c app.Context) {
	status, page := c.Param("status"), c.Param("page")
	email := c.GetString("email")

	pageNO, _ := strconv.Atoi(page)
	if pageNO < 1 {
		c.BadRequest(invalidRequestError)
		return
	}

	if err := ValidateEmail(email); err != nil {
		c.BadRequest(invalidRequestError)
		return
	}

	if !validateCycleStatus(status) {
		c.BadRequest(invalidRequestError)
		return
	}

	res, err := cy.storage.GetAllFromReceiverEmail(status, email, pageNO)
	if err != nil {
		c.StoreError(err)
		return
	}
	result := cy.storage.ToDisplayFormatAll(res)

	// fmt.Println("result: ", result[0])

	c.OK(cy.storage.GetUserDetailWithEmail(result))
}

// GetFromUserEmail godoc
//
//	@summary		GetFromUserEmail
//	@description	Get Cycles by from User email (In context)
//	@tags			cycle
//	@id				GetFromUserEmail
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@response		200	{object}	cycle.CycleDisplay	"Cycle retrieved successfully."
//	@response		400	{object}	app.Response		"Invalid request format or data missing."
//	@response		401	{object}	app.Response		"Authorization failed. Please provide a valid token."
//	@response		404	{object}	app.Response		"Cycle not found with the specified ID."
//	@response		500	{object}	app.Response		"An internal server error occurred while processing the request."
//	@router			/cycles/email/user [get]
func (cy *cycleHandler) GetAllFromUserEmail(c app.Context) {
	email := c.GetString("email")
	res, err := cy.storage.GetFromUserEmail(email)
	if err != nil {
		c.StoreError(err)
		return
	}

	c.OK(cy.storage.ToDisplayFormatAll(res))
}

// @summary		Get a cycle by ID
// @description	Retrieves a cycle by its unique identifier.
// @tags			cycle
// @id				GetOneCycleByID
// @security		BearerAuth
// @accept			json
// @produce		json
// @param			id	path		string			true	"Cycle ID"
// @response		200	{object}	cycle.Cycle		"Cycle retrieved successfully."
// @response		400	{object}	app.Response	"Invalid request format or data missing."
// @response		401	{object}	app.Response	"Authorization failed. Please provide a valid token."
// @response		404	{object}	app.Response	"Cycle not found with the specified ID."
// @response		500	{object}	app.Response	"An internal server error occurred while processing the request."
// @router			/cycles/{id} [get]
func (h *cycleHandler) GetOneByID(c app.Context) {
	id := c.Param("id")
	res, err := h.storage.GetNewByID(id)
	if err != nil {
		if err == invalidRequestError {
			c.BadRequest(err)
			return
		}
		if err == cycleNotFoundError {
			c.NotFound(err)
			return
		}

		c.InternalServerError(err)
		return
	}

	toUserDetail, err := h.storage.ToNewUserDetailFormat(res)
	if err != nil {
		c.InternalServerError(err)
		return
	}

	c.OK(toUserDetail)
}

// DeleteByID godoc
//
//	@summary		DeleteByID
//	@description	Delete Cycles by cycles ID
//	@tags			cycle
//	@id				DeleteByID
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@param			squadID	path		string			true	"Squad ID"
//	@response		200		{object}	app.Response	"nil"
//	@response		400		{object}	app.Response	"invalid cycle id"
//	@response		404		{object}	app.Response	"cycle not found"
//	@response		500		{object}	app.Response	"Internal Server Error"
//	@router			/cycles/{id} [delete]
func (h *cycleHandler) DeleteByID(c app.Context) {
	id := c.Param("id")
	err := h.storage.DeleteOne(id)
	if err != nil {
		if err == invalidRequestError {
			c.BadRequest(err)
			return
		}
		if err == cycleNotFoundError {
			c.NotFound(err)
			return
		}

		c.InternalServerError(err)
		return
	}

	c.OK(nil)
}

// GetCycleProgess godoc
//
//	@summary		GetCycleProgess
//	@description	Get Cycle with status On progress
//	@tags			cycle
//	@id				GetCycleProgess
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@param			cycleId	path		string			true	"Cycle Id"
//	@response		200		{object}	app.Response	"nil"
//	@response		400		{object}	app.Response	"invalid cycle id"
//	@response		404		{object}	app.Response	"cycle not found"
//	@response		500		{object}	app.Response	"Internal Server Error"
//	@router			/cycles/progress/{id} [get]
func (u *cycleHandler) GetCycleProgess(c app.Context) {
	id := c.Param("id")
	cycle, err := u.storage.GetByID(id)
	if err != nil {
		c.InternalServerError(err)
		return
	}

	cycleReponse, err := u.storage.ToDisplayFormat(cycle)
	if err != nil {
		c.InternalServerError(err)
		return
	}

	endDate := cycleReponse.EndDate
	var cycleProgress []GetCycleProgress
	for _, cycle := range cycleReponse.QuantitativeSkill {
		cycleProgress = append(cycleProgress, GetCycleProgress{
			Name:          cycle.Name,
			PersonalScore: cycle.PersonalScore,
			GoalScore:     cycle.LeadGoalScore,
			EndDate:       endDate.Format(time.RFC3339),
		})
	}
	c.OK(cycleProgress)
}

// UpdateUserFinalScore godoc
//
//	@summary		UpdateUserFinalScore
//	@description	Update user score by Score in cycle.
//	@tags			cycle
//	@id				UpdateUserFinalScore
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@param			cycleId	path		string			true	"CycleId"
//	@response		200		{object}	app.Response	"nil"
//	@response		400		{object}	app.Response	"invalid cycle id"
//	@response		404		{object}	app.Response	"cycle not found"
//	@response		500		{object}	app.Response	"Internal Server Error"
//	@router			/cycles/update/{id} [post]
func (cy *cycleHandler) UpdateUserFinalScore(c app.Context) {
	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		c.BadRequest(err)
		return
	}

	err := cy.storage.UpdateUserFinalScore(id)
	if err != nil {
		c.StoreError(err)
		return
	}

	c.OK("successfully update")
}

// UpdateHardSkillsByEmail godoc
//
//	@Summary		Update Hard Skills for a User's Active Cycle
//	@Description	Inserts Hard Skills into the existing cycle for the specified user.
//	@Tags			cycle
//	@ID				UpdateHardSkillsByEmail
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			reqJson	body		UpdateGoalSkillsRequest	true	"Hard Skills input"
//	@Success		200		{object}	app.Response			"Successful operation"
//	@Failure		400		{object}	app.Response			"Error marshaling JSON or mismatched goal score"
//	@Failure		404		{object}	app.Response			"Cycle not found"
//	@Failure		500		{object}	app.Response			"Internal Server Error"
//	@Router			/cycles/goal [put]
func (cy *cycleHandler) UpdateHardSkillsByEmail(c app.Context) {
	email := c.GetString("email")

	var skillJson UpdateGoalSkillsRequest
	err := c.ShouldBindJSON(&skillJson)
	if err != nil {
		c.BadRequest(fmt.Errorf("Error marshaling JSON"))
		return
	}

	for _, hardSkill := range skillJson.HardSkills {
		if hardSkill.PersonalScore+1 < hardSkill.GoalScore || hardSkill.PersonalScore > hardSkill.GoalScore {
			c.BadRequest(fmt.Errorf("miss match goal-score"))
			return
		}
	}

	_, err = cy.storage.UpdateHardSkillsByEmail(c.Ctx(), email, skillJson)
	if err != nil {
		c.StoreError(err)
		return
	}
	c.OK(map[string]string{})
}

// GetLatestCycleFromEmail
func (cy *cycleHandler) GetLatestCycleFromUserEmail(c app.Context) {
	email := c.GetString("email")

	cycle, err := cy.storage.GetLatestCycleFromUserEmail(email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.BadRequest(fmt.Errorf("no documents in result"))
			return
		}
		c.BadRequest(fmt.Errorf("unable to get cycle"))
		return
	}
	c.OK(cycle)
}
