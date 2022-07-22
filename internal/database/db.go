package database

import (
	"context"

	"github.com/CayenneLow/codenames-eventrouter/config"
	"github.com/CayenneLow/codenames-eventrouter/internal/event"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database interface {
	Disconnect(ctx context.Context) error
	GetEventsByGameId(ctx context.Context, gameId string) ([]event.Event, error)
	Insert(ctx context.Context, event event.Event) error
}

type database struct {
	dbClient     *mongo.Client
	dbName       string
	dbCollection string
}

func Init(ctx context.Context, cfg config.Config) Database {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DbURI))
	if err != nil {
		log.Fatalf("Unable to connect to DB: %s", err)
	}
	db := database{
		dbClient:     client,
		dbName:       cfg.DbName,
		dbCollection: cfg.DbCollection,
	}

	return &db
}

func (d *database) Disconnect(ctx context.Context) error {
	return d.dbClient.Disconnect(ctx)
}

func (d *database) GetEventsByGameId(ctx context.Context, gameId string) ([]event.Event, error) {
	cur, err := d.getCollection().Find(ctx, bson.D{primitive.E{Key: "_id", Value: gameId}})
	if err != nil {
		return nil, err
	}
	events := []event.Event{}
	err = cur.All(ctx, &events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (d *database) Insert(ctx context.Context, event event.Event) error {
	_, err := d.getCollection().InsertOne(ctx, event.JsonString())
	return err
}

func (d *database) getCollection() *mongo.Collection {
	return d.dbClient.Database(d.dbName).Collection(d.dbCollection)
}
