package profile

import (
	"context"
	"errors"
	"fmt"

	"gitdev.devops.krungthai.com/aster/ariskill/app"
	"gitdev.devops.krungthai.com/aster/ariskill/app/squad"
	"go.mongodb.org/mongo-driver/mongo"
)

type SquadStorage interface {
	GetOneByID(ctx context.Context, id string) (*squad.Squad, error)
	UpdateByID(ctx context.Context, id string, squad *squad.Squad) (*squad.Squad, error)
}

type squadHandler struct {
	storage SquadStorage
}

func NewSquadHandler(st SquadStorage) *squadHandler {
	return &squadHandler{storage: st}
}

// GetUserSkillRatingBySquadID godoc
//
//	@summary		GetUserSkillRatingBySquadID
//	@description	Get one squad with a user skill ratings in that squad
//	@tags			profile
//	@id				GetUserSkillRatingBySquadID
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@param			squadID	path		string						true	"Squad ID"
//	@response		200		{object}	mySkillRateInSquadResponse	"OK"
//	@response		400		{object}	app.Response				"Bad Request"
//	@response		401		{object}	app.Response				"Unauthorized"
//	@response		500		{object}	app.Response				"Internal Server Error"
//	@router			/profile/squad/{squadID}/skill-ratings [get]
func (h squadHandler) GetUserSkillRatingBySquadID(c app.Context) {
	sqid := c.Param("squadID")
	uId := c.GetString("profileID")

	squadInfo, err := h.storage.GetOneByID(c.Ctx(), sqid)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.NotFound(err)
			return
		}
		c.InternalServerError(err)
		return
	}

	data := mySkillRateInSquadResponse{}
	data.ID = squadInfo.Id
	data.Name = squadInfo.Name
	data.Description = squadInfo.Description
	data.CreatedAt = squadInfo.CreatedAt

	var skill Skill
	for i := 0; i < len(squadInfo.SkillsRatings); i++ {
		for j := 0; j < len(squadInfo.SkillsRatings[i].Ratings); j++ {
			if fmt.Sprintf("%v", squadInfo.SkillsRatings[i].Ratings[j].UserId) == uId {
				skill = Skill{
					SkillID: squadInfo.SkillsRatings[i].SkillId,
					Score:   squadInfo.SkillsRatings[i].Ratings[j].Score,
				}
				data.Skills = append(data.Skills, skill)
				break
			}
		}
	}
	c.OK(data)
}

// RateSkills godoc
//
//	@summary		RateSkills
//	@description	Rate every skills
//	@tags			profile
//	@id				RateSkills
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@param			squadID		path		string			true	"Squad ID"
//	@param			rateSkill	body		RateSkill		true	"Rate Skill"
//	@response		200			{object}	squad.Squad		"OK"
//	@response		400			{object}	app.Response	"Bad Request"
//	@response		401			{object}	app.Response	"Unauthorized"
//	@response		500			{object}	app.Response	"Internal Server Error"
//	@router			/profile/squad/{squadID}/skill-ratings [post]
func (h *squadHandler) RateSkills(c app.Context) {
	squadId := c.Param("squadID")
	uId := squad.GoogleUserId(c.GetString("profileID"))
	var rateSkill RateSkill

	err := c.Bind(&rateSkill)
	if err != nil {
		c.BadRequest(err)
		return
	}

	curSq, err := h.storage.GetOneByID(c.Ctx(), squadId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.NotFound(err)
			return
		}
		c.InternalServerError(err)
		return
	}

	*curSq = updateSquadSkillsRatings(uId, rateSkill, *curSq)

	result, err := h.storage.UpdateByID(c.Ctx(), squadId, curSq)
	if err != nil {
		c.InternalServerError(err)
		return
	}

	c.OK(result)
}

func updateSquadSkillsRatings(uId squad.GoogleUserId, rateSkill RateSkill, curSq squad.Squad) squad.Squad {
	for _, v := range rateSkill.Ratings {
		for j, value := range curSq.SkillsRatings {
			if v.SkillId == value.SkillId {
				isNew := true
				for k, rating := range value.Ratings {
					if rating.UserId == uId {
						curSq.SkillsRatings[j].Ratings[k].Score = v.Score
						isNew = false
						break
					}
				}
				if isNew {
					curSq.SkillsRatings[j].Ratings = append(curSq.SkillsRatings[j].Ratings, squad.SkillRating{UserId: uId, Score: v.Score})
				}
				break
			}
		}
	}

	return curSq
}
