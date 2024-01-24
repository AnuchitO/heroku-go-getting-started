package main

import (
	"context"
	"log"
	"os"
	"time"

	"gitdev.devops.krungthai.com/aster/ariskill/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

	db := client.Database("ariskill")
	err = db.Drop(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
