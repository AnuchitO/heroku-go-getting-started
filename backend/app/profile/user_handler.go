package profile

import (
	"errors"

	"gitdev.devops.krungthai.com/aster/ariskill/app"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStorage interface {
	AboutMeUpdate(id string, about aboutme) error
}

type userHandler struct {
	storage UserStorage
}

func NewUserHandler(st UserStorage) *userHandler {
	return &userHandler{storage: st}
}

// UpdateUser godoc
//
//	@summary		UpdateUser
//	@description	Update user profile
//	@tags			profile
//	@id				UpdateUser
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@param			reqJson	body		aboutme			true	"Editable object"
//	@response		200		{object}	nil				"OK"
//	@response		400		{object}	app.Response	"Bad Request"
//	@response		401		{object}	app.Response	"Unauthorized"
//	@response		405		{object}	app.Response	"Store Error"
//	@response		500		{object}	app.Response	"Internal Server Error"
//	@router			/profile [put]
func (h *userHandler) UpdateAboutMe(c app.Context) {
	var about aboutme
	if err := c.Bind(&about); err != nil {
		c.BadRequest(err)
		return
	}

	id := c.GetString("profileID")
	if err := h.storage.AboutMeUpdate(id, about); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.NotFound(err)
			return
		}
		c.InternalServerError(err)
		return
	}
	c.OK(nil)
}
