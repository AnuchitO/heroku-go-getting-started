package profile

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Profile struct {
	ID              string    `json:"sub" bson:"_id"`
	Email           string    `json:"email" bson:"email"`
	EmployeeID      string    `json:"employeeId" bson:"employee_id"`
	FirstName       string    `json:"givenName" bson:"given_name"`
	LastName        string    `json:"familyName" bson:"family_name"`
	JobRole         string    `json:"jobRole" bson:"job_role"`
	MySquads        []Squad   `json:"squadId" bson:"my_squad"`
	CreatedAt       time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt       time.Time `json:"updatedAt" bson:"updated_at"`
	CreatedBy       string    `json:"createdBy" bson:"created_by"`
	UpdatedBy       string    `json:"updatedBy" bson:"updated_by"`
	AboutMe         string    `json:"aboutMe" bson:"about_me"`
	SocialMedia     []string  `json:"socialMedias" bson:"social_medias"`
	Tags            []string  `json:"tags" bson:"tags"`
	SoftSkills      []Skill   `json:"softSkills" bson:"soft_skills"`
	TechnicalSkills []Skill   `json:"technicalSkills" bson:"technical_skills"`
}
type aboutme struct {
	AboutMe     string   `json:"aboutMe" bson:"about_me"`
	SocialMedia []string `json:"socialMedias" bson:"social_medias"`
	Tags        []string `json:"tags" bson:"tags"`
}
type Squad struct {
	SquadID primitive.ObjectID `json:"sqid,omitempty" bson:"sqid,omitempty"`
	Role    string             `json:"role" bson:"role"`
}
type Skill struct {
	SkillID primitive.ObjectID `json:"skid" bson:"skillID"`
	Score   int                `json:"score" bson:"score"`
}
type SkillsByUser struct {
	UserID string           `bson:"_id"`
	Skills []SkillNameScore `bson:"skills"`
}
type SkillNameScore struct {
	SkillInfo SkillInfo `bson:"skill"`
	Score     int       `bson:"score"`
}
type SkillInfo struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Logo        string             `json:"logo" bson:"logo"`
	Kind        string             `json:"kind" bson:"kind"`
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
type SkillQuerySet struct {
	localField    string
	mapInput      string
	mapInputAlias string
	refSkillID    string
	refSkillScore string
}

type mySkillRateInSquadResponse struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"desc" bson:"desc"`
	CreatedAt   time.Time          `json:"createdAt" bson:"created_at"`
	Skills      []Skill            `json:"skills" bson:"skills"`
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
func NewSkillResponse(s SkillInfo) SkillResponse {
	return SkillResponse{
		ID:          s.ID.Hex(),
		Name:        s.Name,
		Description: s.Description,
		Logo:        s.Logo,
		Kind:        s.Kind,
	}
}

var ErrUserNotFound = errors.New("User not found")
var ErrInvalidKindOfSkill = errors.New("This kind of skill does not exist.")
