package profile

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gitdev.devops.krungthai.com/aster/ariskill/app"
	"gitdev.devops.krungthai.com/aster/ariskill/errs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SkillStorage interface {
	GetSkills(ctx context.Context, id string, kind string) (*SkillsByUser, error)
	UpdateProfileSkills(id string, set skillset, skills []Skill) error
}

type skillHandler struct {
	storage SkillStorage
}

func NewSkillHandler(st SkillStorage) *skillHandler {
	return &skillHandler{storage: st}
}

type skillset string

type RateSkill struct {
	UserId  string `json:"uid" bson:"uid"`
	Ratings []Rating
}

type Rating struct {
	SkillId primitive.ObjectID `json:"skid" bson:"skid"`
	Score   int                `json:"score"`
}

const (
	technical skillset = "technical_skills"
	soft      skillset = "soft_skills"
)

// GetSkillsByUserID godoc
//
//	@summary		GetSkillsByUserID
//	@description	Get skill by user id
//	@tags			profile
//	@id				GetSkillsByUserID
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@param			kind	query		string							false	"Kind of skill [soft , technical]"
//	@response		200		{object}	user.GetSkillByUserIDResponse	"OK"
//	@response		400		{object}	app.Response					"Bad Request"
//	@response		401		{object}	app.Response					"Unauthorized"
//	@response		405		{object}	app.Response					"Store Error"
//	@response		500		{object}	app.Response					"Internal Server Error"
//	@router			/profile/skills [get]
func (h *skillHandler) GetSkillsByUserID(c app.Context) {
	id := c.GetString("profileID")
	kind := c.Query("kind")

	data, err := h.storage.GetSkills(c.Ctx(), id, kind)
	if err != nil {
		if errors.Is(err, ErrInvalidKindOfSkill) {
			c.BadRequest(errs.NewBadRequestError(err.Error()))
			return
		}
		c.InternalServerError(err)
		return
	}
	c.OK(NewSkillByUserResponse(data))
}

// UpdateTechnicalSkills godoc
//
//	@summary		UpdateTechnicalSkills
//	@description	Update user technical skills
//	@tags			profile
//	@id				UpdateTechnicalSkills
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@param			req	body		[]Skill			true	"List of user technical skills"
//	@response		200	{object}	string			"OK"
//	@response		400	{object}	app.Response	"Bad Request"
//	@response		401	{object}	app.Response	"Unauthorized"
//	@response		404	{object}	app.Response	"Not Found"
//	@response		405	{object}	app.Response	"Store Error"
//	@response		500	{object}	app.Response	"Internal Server Error"
//	@router			/profile/skills/technical [post]
func (s *skillHandler) UpdateTechnicalSkill(c app.Context) {
	s.Update(c, technical)
}

// UpdateSoftSkills godoc
//
//	@summary		UpdateSoftSkills
//	@description	Update user soft skills
//	@tags			profile
//	@id				UpdateSoftSkills
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@param			req	body		[]Skill			true	"List of user soft skills"
//	@response		200	{object}	string			"OK"
//	@response		400	{object}	app.Response	"Bad Request"
//	@response		401	{object}	app.Response	"Unauthorized"
//	@response		404	{object}	app.Response	"Not Found"
//	@response		405	{object}	app.Response	"Store Error"
//	@response		500	{object}	app.Response	"Internal Server Error"
//	@router			/profile/skills/soft [post]
func (h *skillHandler) UpdateSoftSkill(c app.Context) {
	h.Update(c, soft)
}
func (s *skillHandler) Update(c app.Context, set skillset) {
	var skills []Skill
	if err := c.Bind(&skills); err != nil {
		c.BadRequest(err)
		return
	}

	id := c.GetString("profileID")

	if err := s.storage.UpdateProfileSkills(id, set, skills); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.NotFound(err)
			return
		}
		c.InternalServerError(err)
		return
	}

	c.OK(fmt.Sprintf("updated %s skill", strings.ReplaceAll(string(set), "_skills", "")))
}
