package skill

import (
	"context"
	"errors"

	"gitdev.devops.krungthai.com/aster/ariskill/app"
	"go.mongodb.org/mongo-driver/mongo"
)

type Storage interface {
	GetByKind(ctx context.Context, kind string) ([]Skill, error)
	GetByID(ctx context.Context, oid string) (Skill, error)
	GetByRole(ctx context.Context, role string) ([]HardSkill, error)
}
type skillHandler struct {
	storage Storage
}

func NewSkillHandler(st Storage) *skillHandler {
	return &skillHandler{
		storage: st,
	}
}

// GetSkillsByKind godoc
//
//	@summary		GetSkillsByKind
//	@description	Get all skill or all of that kind of skill
//	@tags			skill
//	@id				GetSkillsByKind
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@param			kind	query		string			false	"Kind of skill [soft , technical]"
//	@response		200		{object}	skill.Skill		"OK"
//	@response		400		{object}	app.Response	"Bad Request"
//	@response		401		{object}	app.Response	"Unauthorized"
//	@response		405		{object}	app.Response	"Store Error"
//	@response		500		{object}	app.Response	"Internal Server Error"
//	@Router			/skills [get]
func (s *skillHandler) GetSkillsByKind(c app.Context) {
	kind := c.Param("kindtype")
	sks, err := s.storage.GetByKind(c.Ctx(), kind)
	if err != nil {
		c.InternalServerError(err)
		return
	}
	c.OK(sks)
}

// GetSkillByID godoc
//
//	@summary		GetSkillByID
//	@description	Get a skill by skill id
//	@tags			skill
//	@id				GetSkillByID
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@param			id	path		string			true	"Skill ID"
//	@response		200	{object}	skill.Skill		"OK"
//	@response		400	{object}	app.Response	"Bad Request"
//	@response		401	{object}	app.Response	"Unauthorized"
//	@response		405	{object}	app.Response	"Store Error"
//	@response		500	{object}	app.Response	"Internal Server Error"
//	@Router			/skills/{id} [get]
func (h *skillHandler) SkillByID(c app.Context) {
	id := c.Param("id")
	sk, err := h.storage.GetByID(c.Ctx(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.NotFound(err)
			return
		}
		c.InternalServerError(err)
		return
	}

	c.OK(sk)
}

// SkillByJobRole find all skills by jobrole of transaction user
func (h *skillHandler) SkillByJobRole(c app.Context) {
	role := c.GetString("role")
	sk, err := h.storage.GetByRole(c.Ctx(), role)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.NotFound(err)
			return
		}
		c.InternalServerError(err)
		return
	}

	c.OK(sk)
}
