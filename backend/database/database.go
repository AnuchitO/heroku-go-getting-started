package database

import (
	"context"
	"fmt"
	"time"

	"gitdev.devops.krungthai.com/aster/ariskill/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongo(cfg config.Database) (*mongo.Database, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	auth := options.Credential{
		Username: cfg.Username,
		Password: cfg.Password,
	}
	option := options.Client().
		ApplyURI(cfg.Host).
		SetAuth(auth).
		SetConnectTimeout(10 * time.Second)

	conn, err := mongo.Connect(ctx, option)
	if err != nil {
		panic(fmt.Errorf("can not connect to mongodb: %w", err))
	}

	teardown := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = conn.Disconnect(ctx)
	}

	return conn.Database(cfg.Name), teardown
}
