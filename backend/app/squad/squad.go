package squad

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GoogleUserId string

type Squad struct {
	Id            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name" bson:"name"`
	TeamleadMail  string             `json:"teamleadMail" bson:"teamleadMail"`
	Description   string             `json:"desc" bson:"desc"`
	CreatedAt     time.Time          `json:"createdAt" bson:"created_at"`
	SkillsRatings []SkillRatings     `json:"skillsRatings" bson:"skills_ratings"`
}

type SquadFilter struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"desc" bson:"desc"`
}

type SkillRatings struct {
	SkillId primitive.ObjectID `json:"skid" bson:"skid"`
	Ratings []SkillRating      `json:"ratings" bson:"ratings"`
}

type SkillRating struct {
	UserId GoogleUserId `json:"uid" bson:"uid"`
	Score  int          `json:"score" bson:"score"`
}

type SquadSkill struct {
	SkillId primitive.ObjectID `json:"skid" bson:"skid"`
	Ratings []SkillRating      `json:"ratings" bson:"ratings"`
}

// ========= Input =========
// Input body of user rating a Squad skill
type SkillRatingInput struct {
	UserId  primitive.ObjectID `json:"uid" bson:"uid"`
	SkillId primitive.ObjectID `json:"skid" bson:"skid"`
	Score   int                `json:"score" bson:"score"`
}

// ======== Output =========
// Output body of Squad with averages skills' rating
type SquadAverageSkillOutput struct {
	SquadId       primitive.ObjectID `json:"squadId"`
	AveragesSkill []AverageSkill     `json:"averagesSkill"`
}

type AverageSkill struct {
	SkillId primitive.ObjectID `json:"skid"`
	Average float64            `json:"average"`
}

type TeamLeadDetail struct {
	Name string
	Mail string
}
