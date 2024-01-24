package cycle

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	StatusAll      = "All"
	StatusPending  = "Pending"
	StatusApproved = "Approved"
	StatusRunning  = "In Progress"
	StatusDone     = "Done"
)

const (
	StateReview = "Review"
)

type Cycle struct {
	ID                primitive.ObjectID  `json:"id,omitempty" bson:"_id,omitempty"`
	SenderMail        string              `json:"senderMail" bson:"sender_mail" binding:"required"`
	ReceiverMail      string              `json:"receiverMail" bson:"receiver_mail" binding:"required"`
	StartDate         time.Time           `json:"startDate" bson:"start_date" binding:"required"`
	EndDate           time.Time           `json:"endDate" bson:"end_date" binding:"required"`
	QuantitativeSkill []QuantitativeSkill `json:"quantitativeSkill" bson:"quantitative_skill"`
	IntuitiveSkill    []IntuitiveSkill    `json:"intuitiveSkill" bson:"intuitive_skill"`
	Status            string              `json:"status" bson:"status" binding:"required"`
	State             string              `json:"state" bson:"state"`
	Comment           string              `json:"comment" bson:"comment"`
}

type CycleSave struct {
	ID                primitive.ObjectID  `json:"id,omitempty" bson:"_id,omitempty"`
	SenderMail        string              `json:"senderMail" bson:"sender_mail" binding:"required"`
	ReceiverMail      string              `json:"receiverMail" bson:"receiver_mail" binding:"required"`
	StartDate         time.Time           `json:"startDate" bson:"start_date" binding:"required"`
	EndDate           time.Time           `json:"endDate" bson:"end_date" binding:"required"`
	QuantitativeSkill []QuantitativeSkill `json:"quantitativeSkill" bson:"quantitative_skill"`
	IntuitiveSkill    []IntuitiveSkill    `json:"intuitiveSkill" bson:"intuitive_skill"`
	Comment           string              `json:"comment" bson:"comment"`
}

type CycleDisplay struct {
	ID                primitive.ObjectID         `json:"id" bson:"_id"`
	SenderMail        string                     `json:"senderMail" bson:"sender_mail"`
	ReceiverMail      string                     `json:"receiverMail" bson:"receiver_mail"`
	StartDate         time.Time                  `json:"startDate" bson:"start_date"`
	EndDate           time.Time                  `json:"endDate" bson:"end_date"`
	QuantitativeSkill []QuantitativeSkillDisplay `json:"quantitativeSkill" bson:"quantitative_skill"`
	IntuitiveSkill    []IntuitiveSkill           `json:"intuitiveSkill" bson:"intuitive_skill"`
	Status            string                     `json:"status" bson:"status"`
	State             string                     `json:"state" bson:"state"`
	Comment           string                     `json:"comment" bson:"comment"`
}

type QuantitativeSkill struct {
	ID            primitive.ObjectID `json:"id" bson:"id"`
	PersonalScore int                `json:"personalScore" bson:"personal_score" binding:"required"`
	GoalScore     int                `json:"goalScore" bson:"goal_score" binding:"required"`
	LeadGoalScore int                `json:"leadGoalScore" bson:"lead_goal_score"`
	FinalScore    int                `json:"finalScore" bson:"final_score"`
	Comment       string             `json:"comment" bson:"comment"`
}

type QuantitativeSkillDisplay struct {
	ID            primitive.ObjectID `json:"id" bson:"id"`
	Name          string             `json:"name" bson:"name"`
	Description   string             `json:"description" bson:"description"`
	Logo          string             `json:"logo" bson:"logo"`
	PersonalScore int                `json:"personalScore" bson:"personal_score" binding:"required"`
	GoalScore     int                `json:"goalScore" bson:"goal_score" binding:"required"`
	LeadGoalScore int                `json:"leadGoalScore" bson:"lead_goal_score"`
	FinalScore    int                `json:"finalScore" bson:"final_score"`
	Comment       string             `json:"comment" bson:"comment"`
}

type CycleInput struct {
	ReceiverMail string    `json:"receiverMail" bson:"receiver_mail" binding:"required"`
	StartDate    time.Time `json:"startDate" bson:"start_date" binding:"required"`
	EndDate      time.Time `json:"endDate" bson:"end_date" binding:"required"`
	// QuantitiveSkill >= 1 skill
	QuantitativeSkill []QuantitativeSkill `json:"quantitativeSkill" bson:"quantitative_skill" binding:"required"`
	IntuitiveSkill    []IntuitiveSkill    `json:"intuitiveSkill" bson:"intuitive_skill"`
	Comment           string              `json:"comment" bson:"comment"`
}

type IntuitiveSkill struct {
	Name    string `json:"name" bson:"name" binding:"required"`
	Status  string `json:"status" bson:"status"`
	Goal    string `json:"goal" bson:"goal" binding:"required"`
	Comment string `json:"comment" bson:"comment"`
}

type QuantiativeSkillToUpdate struct {
	ID         primitive.ObjectID `json:"id" bson:"id"`
	FinalScore int                `json:"finalScore" bson:"final_score"`
}

type GetFinalScoreResponse struct {
	Email             string                     `json:"email" binding:"required"`
	QuantiativeSkills []QuantiativeSkillToUpdate `json:"quantiativeSkills"`
}

type CycleWithUserDetail struct {
	ID                primitive.ObjectID         `json:"id" bson:"_id"`
	UID               string                     `json:"sub" bson:"user_id"`
	FirstName         string                     `json:"givenName" bson:"given_name"`
	LastName          string                     `json:"familyName" bson:"family_name"`
	JobRole           string                     `json:"jobRole" bson:"job_role"`
	Level             string                     `json:"level" bson:"level"`
	SenderMail        string                     `json:"senderMail" bson:"sender_mail"`
	ReceiverMail      string                     `json:"receiverMail" bson:"receiver_mail"`
	StartDate         time.Time                  `json:"startDate" bson:"start_date"`
	EndDate           time.Time                  `json:"endDate" bson:"end_date"`
	QuantitativeSkill []QuantitativeSkillDisplay `json:"quantitativeSkill" bson:"quantitative_skill"`
	IntuitiveSkill    []IntuitiveSkill           `json:"intuitiveSkill" bson:"intuitive_skill"`
	Status            string                     `json:"status" bson:"status"`
	State             string                     `json:"state" bson:"state"`
	Comment           string                     `json:"comment" bson:"comment"`
}

// new cycle
type NewCycle struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TeamLeaderMail string             `json:"teamLeaderMail" bson:"teamLeaderMail" binding:"required"`
	AriserMail     string             `json:"ariserMail" bson:"ariserMail" binding:"required"`
	StartDate      time.Time          `json:"startDate" bson:"startDate" binding:"required"`
	EndDate        time.Time          `json:"endDate" bson:"endDate" binding:"required"`
	Status         string             `json:"status" bson:"status" binding:"required"`
	Comment        string             `json:"comment" bson:"comment"`
	HardSkills     []HardSkill        `json:"hardSkills" bson:"hardSkills,omitempty"`
	State          string             `json:"state" bson:"state"`
}

type NewCycleDisplay struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TeamLeaderMail string             `json:"teamLeaderMail" bson:"teamLeaderMail"`
	AriserMail     string             `json:"ariserMail" bson:"ariserMail" `
	StartDate      time.Time          `json:"startDate" bson:"startDate" `
	EndDate        time.Time          `json:"endDate" bson:"endDate"`
	Status         string             `json:"status" bson:"status"`
	Comment        string             `json:"comment" bson:"comment"`
	HardSkills     []HardSkillDisplay `json:"hardSkills,omitempty" bson:"hardSkills,omitempty"`
	State          string             `json:"state" bson:"state"`
}

type NewCycleWithUserDetail struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FirstName      string             `json:"givenName" bson:"giveName"`
	LastName       string             `json:"familyName" bson:"familyName"`
	JobRole        string             `json:"jobRole" bson:"jobRole"`
	Level          string             `json:"level" bson:"level"`
	TeamLeaderMail string             `json:"teamLeaderMail" bson:"teamLeaderMail"`
	AriserMail     string             `json:"ariserMail" bson:"ariserMail" `
	StartDate      time.Time          `json:"startDate" bson:"startDate" `
	EndDate        time.Time          `json:"endDate" bson:"endDate"`
	Status         string             `json:"status" bson:"status"`
	Comment        string             `json:"comment" bson:"comment"`
	HardSkills     []HardSkillDisplay `json:"hardSkills,omitempty" bson:"hardSkills,omitempty"`
	State          string             `json:"state" bson:"state"`
}

type UpdateGoalSkillsRequest struct {
	HardSkills []HardSkill `json:"hardSkills" bson:"hardSkills" binding:"required"`
}

type HardSkill struct {
	ID            primitive.ObjectID `json:"id" bson:"id"`
	Name          string             `json:"name" bson:"name" binding:"required"`
	Description   string             `json:"description" bson:"description" binding:"required"`
	SkillLevels   []SkillLevel       `json:"skillLevels" bson:"skillLevels" binding:"required"`
	PersonalScore int                `json:"personalScore" bson:"personalScore" binding:"required"`
	GoalScore     int                `json:"goalScore" bson:"goalScore" binding:"required"`
	LeadScore     int                `json:"leadScore" bson:"leadScore"`
	MutualScore   int                `json:"mutualScore" bson:"mutualScore"`
	Comment       string             `json:"comment" bson:"comment"`
}
type HardSkillDisplay struct {
	ID            primitive.ObjectID `json:"id" bson:"id"`
	Name          string             `json:"name" bson:"name" `
	Description   string             `json:"description" bson:"description" `
	SkillLevels   []SkillLevel       `json:"skillLevels" bson:"skillLevels" `
	PersonalScore int                `json:"personalScore" bson:"personalScore" `
	GoalScore     int                `json:"goalScore" bson:"goalScore" `
}

type SkillLevel struct {
	Level            int    `json:"level" bson:"level" binding:"required"`
	LevelDescription string `json:"levelDescription" bson:"levelDescription" binding:"required"`
}
