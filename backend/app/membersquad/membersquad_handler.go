package membersquad

import (
	"context"
	"fmt"

	"gitdev.devops.krungthai.com/aster/ariskill/app"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Storage interface {
	GetByID(ctx context.Context, id string) (*Member, error)
	UpdateMySquad(ctx context.Context, id string, update []Squad) error
	GetAllBySquadId(ctx context.Context, squadId primitive.ObjectID) ([]Member, error)
}

type memberSquadHandler struct {
	storage Storage
}

func NewMemberSquadHandler(st Storage) *memberSquadHandler {
	return &memberSquadHandler{
		storage: st,
	}
}

// AddSquad godoc
//
//	@summary		AddMemberSquad
//	@description	Add member to squad
//	@tags			membersquad
//	@id				AddMemberSquad
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@param			s	body		SquadMember		true	"SquadMember Object"
//	@response		200	{object}	nil				"OK"
//	@response		400	{object}	app.Response	"Bad Request"
//	@response		401	{object}	app.Response	"Unauthorized"
//	@response		500	{object}	app.Response	"Internal Server Error"
//	@router			/member-squads/members [put]
func (h *memberSquadHandler) AddMemberSquad(ctx app.Context) {
	var s SquadMember
	if err := ctx.Bind(&s); err != nil {
		ctx.BadRequest(ErrRequestInvalidFormat)
		return
	}

	if s.IncludeMyself {
		profileId := ctx.GetString("profileID")
		s.Members = append([]UserAndRole{{UserId: profileId, Role: "Leader"}}, s.Members...)
	}

	for _, v := range s.Members {
		u, err := h.storage.GetByID(ctx.Ctx(), v.UserId)
		if err != nil {
			ctx.InternalServerError(err)
			return
		}

		isExists := false
		for _, sq := range u.MySquads {
			if sq.SquadID == s.SquadId {
				isExists = true
			}
		}
		if isExists {
			continue
		}

		u.MySquads = append(u.MySquads, Squad{SquadID: s.SquadId, Role: v.Role})

		err = h.storage.UpdateMySquad(ctx.Ctx(), u.ID, u.MySquads)
		if err != nil {
			ctx.InternalServerError(err)
			return
		}
	}

	ctx.OK(nil)
}

// DeleteMemberSquad godoc
//
//	@summary		DeleteMemberSquad
//	@description	Delete all member from squad (at least 1 member in squad)
//	@tags			membersquad
//	@id				DeleteMemberSquad
//	@security		BearerAuth
//	@accept			json
//	@produce		json
//	@param			squadID	path		string			true	"Squad ID"
//	@response		200		{object}	nil				"OK"
//	@response		400		{object}	app.Response	"Bad Request"
//	@response		401		{object}	app.Response	"Unauthorized"
//	@response		500		{object}	app.Response	"Internal Server Error"
//	@router			/member-squads/{squadID}/members [delete]
func (h *memberSquadHandler) DeleteMemberSquad(actx app.Context) {
	sqId := actx.Param("squadID")
	fmt.Println(sqId)
	objId, err := primitive.ObjectIDFromHex(sqId)
	if err != nil {
		actx.InternalServerError(err)
		return
	}
	users, errRes := h.storage.GetAllBySquadId(actx.Ctx(), objId)
	if errRes != nil {
		actx.InternalServerError(errRes)
		return
	}

	// For Loop to remove squad from User MySquad
	for i, user := range users {
		for j, squad := range user.MySquads {
			if squad.SquadID == objId {
				users[i].MySquads = append(users[i].MySquads[:j], users[i].MySquads[j+1:]...)
				err = h.storage.UpdateMySquad(actx.Ctx(), user.ID, users[i].MySquads)
				if err != nil {
					actx.InternalServerError(err)
					return
				}
				break
			}
		}
	}
	actx.OK(nil)
}
