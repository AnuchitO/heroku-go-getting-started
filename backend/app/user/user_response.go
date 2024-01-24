package user

import (
	"time"

	"gitdev.devops.krungthai.com/aster/ariskill/app/skill"
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

type GetSkillByUserIDResponse struct {
	UserID string                   `json:"id"`
	Skills []SkillNameScoreResponse `json:"skills"`
}

type SkillNameScoreResponse struct {
	SkillInfo SkillResponse `json:"skill"`
	Score     int           `json:"score"`
}

type SkillResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	Kind        string `json:"kind"`
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

func NewSkillByUserResponse(s *SkillsByUser) GetSkillByUserIDResponse {
	return GetSkillByUserIDResponse{
		UserID: s.UserID,
		Skills: NewSkillNameScoreResponse(s.Skills),
	}
}

func NewSkillNameScoreResponse(skill []SkillNameScore) []SkillNameScoreResponse {
	var skills []SkillNameScoreResponse
	for _, s := range skill {
		skills = append(skills, SkillNameScoreResponse{
			Score:     s.Score,
			SkillInfo: NewSkillResponse(s.SkillInfo),
		})
	}
	return skills
}

func NewSkillResponse(s skill.Skill) SkillResponse {
	return SkillResponse{
		ID:          s.ID.Hex(),
		Name:        s.Name,
		Description: s.Description,
		Logo:        s.Logo,
	}
}
