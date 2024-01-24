package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"gitdev.devops.krungthai.com/aster/ariskill/app/cycle"
	"gitdev.devops.krungthai.com/aster/ariskill/app/skill"
	"gitdev.devops.krungthai.com/aster/ariskill/app/squad"
	"gitdev.devops.krungthai.com/aster/ariskill/app/user"
	"gitdev.devops.krungthai.com/aster/ariskill/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// struct for mock employee data
type Employee struct {
	Email      string               `json:"email" bson:"email"`
	EmployeeID string               `json:"employeeId" bson:"employee_id"`
	JobRole    string               `json:"jobRole" bson:"job_role"`
	Project    string               `json:"project" bson:"project"`
	Team       string               `json:"team" bson:"team"`
	SquadID    []primitive.ObjectID `json:"squadId"`
}

func main() {
	cfg := config.C(os.Getenv("ENV"))
	host := cfg.Database.Host
	username := cfg.Database.Username
	password := cfg.Database.Password

	auth := options.Credential{
		Username: username,
		Password: password,
	}
	option := options.Client().
		ApplyURI(host).
		SetAuth(auth).
		SetConnectTimeout(10 * time.Second)

	client, err := mongo.Connect(context.Background(), option)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := client.Disconnect(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}()

	db := client.Database("ariskill")

	usersData := readFile(pathFile("./ariskill.users.json"))
	seed[user.User](usersData, db.Collection("users"))

	skillsData := readFile(pathFile("./ariskill.skills.json"))
	seed[skill.Skill](skillsData, db.Collection("skills"))

	hardSkillsData := readFile(pathFile("./ariskill.skills.hard.json"))
	seed[skill.HardSkill](hardSkillsData, db.Collection("hard_skills"))

	// seed for mock employee data that should be get from company
	// but latest version in batch2 we remove that feature.
	// CenterData := readFile(pathFile("./ariskill.employee.json"))
	// seed[Employee](CenterData, db.Collection("employee"))

	squadData := readFile(pathFile("./ariskill.squads.json"))
	seed[squad.Squad](squadData, db.Collection("squads"))

	cycleData := readFile(pathFile("./ariskill.cycles.json"))
	seed[cycle.Cycle](cycleData, db.Collection("cycles"))

	newCycleData := readFile(pathFile("./ariskill.newcycles.json"))
	seed[cycle.NewCycle](newCycleData, db.Collection("new_cycles"))
}

func seed[T any](data []byte, coln *mongo.Collection) {
	var rows []T
	if err := json.Unmarshal(data, &rows); err != nil {
		log.Fatal(err)
	}

	ids := make([]interface{}, len(rows))
	for _, row := range rows {
		id, err := coln.InsertOne(context.Background(), row)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				continue
			}
			log.Printf("Error inserting data into collection: %v", err)
		}
		ids = append(ids, id)
	}

	fmt.Printf("Inserted %v %v documents\n", coln.Name(), len(ids))
}

func readFile(filename string) []byte {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func pathFile(filename string) string {
	_, fname, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Error getting the current file's path")
		return ""
	}
	return filepath.Join(filepath.Dir(fname), filename)
}
