package profile

import (
	"time"
)

type GetUserResponse struct {
	ID             string        `json:"sub"`
	Email          string        `json:"email"`
	EmployeeID     string        `json:"employeeId"`
	FirstName      string        `json:"givenName"`
	LastName       string        `json:"familyName"`
	JobRole        string        `json:"jobRole"`
	Level          string        `json:"level"`
	AboutMe        string        `json:"aboutMe"`
	MySquad        []MySquad     `json:"squadId"`
	SocialMedia    []string      `json:"socialMedias"`
	Tags           []string      `json:"tags"`
	CreatedAt      time.Time     `json:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt"`
	CreatedBy      string        `json:"createdBy"`
	UpdatedBy      string        `json:"updatedBy"`
	SoftSkill      []MySkill     `json:"softSkills"`
	TechnicalSkill []MySkill     `json:"technicalSkills"`
	HardSkills     []MyHardSkill `json:"hardSkills"`
}

type GetEmailNameResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type MyHardSkillResponse struct {
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	CurrentLevel int          `json:"currentLevel"`
	SkillLevel   []SkillLevel `json:"skillLevel"`
	Sort         int          `json:"sort"`
}

type SkillLevelResponse struct {
	Level            int    `json:"level"`
	LevelDescription string `json:"levelDescription"`
}

func NewGetUserResponse(u User) GetUserResponse {
	return GetUserResponse(u)
}
