package user

import (
	"errors"
	"time"

	"gitdev.devops.krungthai.com/aster/ariskill/app/skill"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             string        `json:"sub" bson:"_id"`
	Email          string        `json:"email" bson:"email"`
	EmployeeID     string        `json:"employeeId" bson:"employee_id"`
	FirstName      string        `json:"givenName,omitempty" bson:"given_name,omitempty"`
	LastName       string        `json:"familyName,omitempty" bson:"family_name,omitempty"`
	JobRole        string        `json:"jobRole,omitempty" bson:"job_role,omitempty"`
	Level          string        `json:"level" bson:"level"`
	AboutMe        string        `json:"aboutMe,omitempty" bson:"about_me,omitempty"`
	MySquad        []MySquad     `json:"mySquad,omitempty" bson:"my_squad,omitempty"`
	SocialMedia    []string      `json:"socialMedias,omitempty" bson:"social_medias,omitempty"`
	Tags           []string      `json:"tags,omitempty" bson:"tags,omitempty"`
	CreatedAt      time.Time     `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt      time.Time     `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
	CreatedBy      string        `json:"createdBy,omitempty" bson:"created_by,omitempty"`
	UpdatedBy      string        `json:"updatedBy,omitempty" bson:"updated_by,omitempty"`
	SoftSkill      []MySkill     `json:"softSkills,omitempty" bson:"soft_skills,omitempty"`
	TechnicalSkill []MySkill     `json:"technicalSkills,omitempty" bson:"technical_skills,omitempty"`
	HardSkills     []MyHardSkill `json:"hardSkills" bson:"hard_skills,omitempty"`
}

type Employee struct {
	Email      string               `json:"email" bson:"email"`
	EmployeeID string               `json:"employeeId" bson:"employee_id"`
	JobRole    string               `json:"jobRole" bson:"job_role"`
	Project    string               `json:"project" bson:"project"`
	Team       string               `json:"team" bson:"team"`
	SquadID    []primitive.ObjectID `json:"squadId" bson:"squadid"`
}

type MySkill struct {
	SkillID primitive.ObjectID `json:"skid" bson:"skillID"`
	Score   int                `json:"score" bson:"score"`
}

type MyHardSkill struct {
	Name         string       `json:"name" bson:"name"`
	Description  string       `json:"description" bson:"description"`
	CurrentLevel int          `json:"currentLevel" bson:"currentLevel"`
	SkillLevel   []SkillLevel `json:"skilllevel" bson:"skilllevel"`
	Sort         int          `json:"sort" bson:"sort"`
}

type SkillLevel struct {
	Level            int    `json:"level" bson:"level"`
	LevelDescription string `json:"leveldescription" bson:"leveldescription"`
}

type SkillsByUser struct {
	UserID string           `bson:"_id"`
	Skills []SkillNameScore `bson:"skills"`
}

type SkillNameScore struct {
	SkillInfo skill.Skill `bson:"skill"`
	Score     int         `bson:"score"`
}

type MySquad struct {
	SquadID primitive.ObjectID `json:"sqid,omitempty" bson:"sqid,omitempty"`
	Role    string             `json:"role" bson:"role"`
}

var ErrInvalidKindOfSkill = errors.New("this kind of skill does not exist")
