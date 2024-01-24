package user

import (
	"context"

	"net/http"
	"strings"

	"gitdev.devops.krungthai.com/aster/ariskill/app"
)

type UserStorage interface {
	GetAll(ctx context.Context) ([]User, error)
	GetOneById(ctx context.Context, id string) (*User, error)
	GetHardSkillById(ctx context.Context, Name string) ([]User, error)
	GetOneByEmail(ctx context.Context, email string) (*User, error)
}

type userHandler struct {
	userStorage UserStorage
}

func NewUserHandler(userStorage UserStorage) *userHandler {
	return &userHandler{
		userStorage: userStorage,
	}
}

// GetAllUsers godoc
//
//	@summary		GetAllUsers
//	@description	Get user profile
//	@tags			user
//	@id				GetAllUsers
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@param			search	query		string					false	"Search name or email"
//	@response		200		{array}		user.GetUserResponse	"OK"
//	@response		400		{object}	app.Response			"Bad Request"
//	@response		401		{object}	app.Response			"Unauthorized"
//	@response		405		{object}	app.Response			"Store Error"
//	@response		500		{object}	app.Response			"Internal Server Error"
//	@router			/users [get]
func (u *userHandler) GetAllUsers(c app.Context) {
	allUsers, err := u.userStorage.GetAll(c.Ctx())
	if err != nil {
		c.InternalServerError(err)
		return
	}

	searchQuery := strings.ToLower(c.Query("search"))

	var userResponses []GetUserResponse
	for _, user := range allUsers {
		if strings.Contains(strings.ToLower(user.FirstName), searchQuery) ||
			strings.Contains(strings.ToLower(user.LastName), searchQuery) ||
			strings.Contains(strings.ToLower(user.Email), searchQuery) {
			userResponses = append(userResponses, NewGetUserResponse(user))
		}
	}
	c.OK(userResponses)
}

// GetUserData godoc
//
//	@summary		Retrieve all users' email and name
//	@description	Fetches all users.
//	@tags			user
//	@id				GetUsersData
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@response		200	{object}	[]GetEmailNameResponse	"OK"
//	@response		401	{object}	app.Response			"Unauthorized"
//	@response		500	{object}	app.Response			"Internal Server Error"
//	@router			/users/email [get]
func (u *userHandler) GetUsersData(c app.Context) {
	allUsers, err := u.userStorage.GetAll(c.Ctx())
	if err != nil {
		c.InternalServerError(err)
		return
	}

	var userResponses []GetEmailNameResponse
	for _, user := range allUsers {
		userResponses = append(userResponses, GetEmailNameResponse{
			Email: user.Email,
			Name:  user.FirstName + " " + user.LastName,
		})
	}
	c.OK(userResponses)
}

// GetMySquadThatImMemberByID godoc
//
//	@summary		Retrieve all team name that Im member
//	@description	Fetches all team.
//	@tags			user
//	@id				GetMySquadThatImMemberByID
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@response		200	{object}	[]GetEmailNameResponse	"OK"
//	@response		401	{object}	app.Response			"Unauthorized"
//	@response		500	{object}	app.Response			"Internal Server Error"
//	@router			/users/:id [get]
func (u *userHandler) GetMySquadThatImMemberByID(c app.Context) {
	id := c.Param("id")
	var squads []MySquad
	res, err := u.userStorage.GetOneById(c.Ctx(), id)
	if err != nil {
		c.NotFound(err)
		return
	}

	if len(res.MySquad) < 1 {
		c.NotFound(teamNotFoundError)
		return
	}

	for _, item := range res.MySquad {
		if item.Role == "member" {
			squads = append(squads, MySquad{
				SquadID: item.SquadID,
				Role:    item.Role,
			})
		}
	}
	if len(squads) < 1 {
		c.BadRequest(IsLeadFoundError)
		return
	}

	c.OK(squads)
}

// GetHardSkillsMemberByID godoc
//	@summary		Retrieve all HardSkill
//	@description	Fetches all HardSkill.
//	@tags			user
//	@id				GetHardSkillsMemberByID
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@response		200	{object}	[]MyHardSkill	"OK"
//	@response		404	{object}	app.Response	"Not Found"
//	@response		401	{object}	app.Response	"Unauthorized"
//	@response		400	{object}	app.Response	"Bad Request"
//	@response		500	{object}	app.Response	"Internal Server Error"
//	@router			/users/{id} [get]
func (u *userHandler) GetHardSkillsMemberByID(c app.Context) {
	id := c.GetString("profileID")
	user, err := u.userStorage.GetOneById(c.Ctx(), id)
	if err != nil {
		c.NotFound(err) // You might want to use the appropriate error response depending on your application
		return
	}

	c.JSON(http.StatusOK, user.HardSkills)
}

func (u *userHandler) GetUserByEmailToken(c app.Context) {
	email := c.GetString("email")

	user, err := u.userStorage.GetOneByEmail(c.Ctx(), email)
	if err != nil {
		c.NotFound(err)
	}

	c.OK(user)
}
