package database

import (
	"context"

	"github.com/CayenneLow/codenames-eventrouter/config"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database interface {
	Disconnect(ctx context.Context) error
}

func Init(ctx context.Context, cfg config.Config) Database {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DbURI))
	if err != nil {
		log.Fatalf("Unable to connect to DB: %s", err)
	}

	collection := client.Database("test").Collection("delete_me")
	collection.InsertOne(ctx, bson.D{{"_id", "T35T1"}, {"test", 123}})
	return client
}
