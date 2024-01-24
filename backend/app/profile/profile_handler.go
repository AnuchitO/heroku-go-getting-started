package profile

import (
	"context"

	"gitdev.devops.krungthai.com/aster/ariskill/app"
)

type Storage interface {
	GetByID(ctx context.Context, id string) (*Profile, error)
}

type profileHandler struct {
	storage Storage
}

func NewProfileHandler(st Storage) *profileHandler {
	return &profileHandler{
		storage: st,
	}
}

// GetUserByID godoc
//
//	@summary		GetUserByID
//	@description	Get user by id
//	@tags			profile
//	@id				GetUserByID
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@response		200	{object}	user.GetUserResponse	"OK"
//	@response		400	{object}	app.Response			"Bad Request"
//	@response		401	{object}	app.Response			"Unauthorized"
//	@response		405	{object}	app.Response			"Store Error"
//	@response		500	{object}	app.Response			"Internal Server Error"
//	@router			/profile [get]
func (s *profileHandler) User(c app.Context) {
	id := c.GetString("profileID")

	user, err := s.storage.GetByID(c.Ctx(), id)
	if err != nil {
		c.InternalServerError(err)
		return
	}

	c.OK(user)
}
