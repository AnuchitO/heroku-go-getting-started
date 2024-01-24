package membersquad

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Member struct {
	ID              string    `json:"sub" bson:"_id"`
	Email           string    `json:"email" bson:"email"`
	EmployeeID      string    `json:"employee_id" bson:"employee_id"`
	FirstName       string    `json:"given_name" bson:"given_name"`
	LastName        string    `json:"family_name" bson:"family_name"`
	JobRole         string    `json:"jobRole" bson:"job_role"`
	MySquads        []Squad   `json:"squadId" bson:"my_squad"`
	CreatedAt       time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt       time.Time `json:"updatedAt" bson:"updated_at"`
	CreatedBy       string    `json:"createdBy" bson:"created_by"`
	UpdatedBy       string    `json:"updatedBy" bson:"updated_by"`
	SoftSkills      []Skill   `json:"softSkills" bson:"soft_skills"`
	TechnicalSkills []Skill   `json:"technicalSkills" bson:"technical_skills"`
}

type Squad struct {
	SquadID primitive.ObjectID `json:"sqid,omitempty" bson:"sqid,omitempty"`
	Role    string             `json:"role" bson:"role"`
}

type Skill struct {
	SkillID primitive.ObjectID `json:"skid" bson:"skillID"`
	Score   int                `json:"score" bson:"score"`
}

type UserAndRole struct {
	UserId string `json:"uid"`
	Role   string `json:"role"`
}

type SquadMember struct {
	SquadId       primitive.ObjectID `json:"sqid"`
	Members       []UserAndRole      `json:"members"`
	IncludeMyself bool               `json:"include"`
}

var ErrRequestInvalidFormat = errors.New("Request is invalid format")
