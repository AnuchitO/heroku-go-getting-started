package skill

import "go.mongodb.org/mongo-driver/bson/primitive"

type SkillKind string

type Skill struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Logo        string             `json:"logo" bson:"logo"`
	Kind        string             `json:"kind" bson:"kind"`
}

type HardSkill struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description DescriptionEnum    `json:"description" bson:"description"`
	JobRole     []JobRole          `json:"jobRole" bson:"jobRole"`
	Sort        int                `json:"sort" bson:"sort"`
	SkillLevel  []SkillLevel       `json:"skillLevel" bson:"skillLevel"`
}

type ID struct {
	OID string `json:"$oid"`
}

type SkillLevel struct {
	Level            int              `json:"level"`
	LevelDescription LevelDescription `json:"levelDescription"`
}

type DescriptionEnum string

const (
	Description DescriptionEnum = "description"
)

type JobRole string

const (
	Backend   JobRole = "backend"
	Frontend  JobRole = "frontend"
	Fullstack JobRole = "fullstack"
)

type LevelDescription string

const (
	ExampleLevel1 LevelDescription = "example level 1"
	ExampleLevel2 LevelDescription = "example level 2"
	ExampleLevel3 LevelDescription = "example level 3"
	ExampleLevel4 LevelDescription = "example level 4"
	ExampleLevel5 LevelDescription = "example level 5"
)
