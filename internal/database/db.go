package database

import (
	"context"

	"github.com/CayenneLow/codenames-eventrouter/config"
	"github.com/CayenneLow/codenames-eventrouter/pkg/event"
	redis "github.com/go-redis/redis/v9"
)

type Database interface {
	Disconnect(ctx context.Context) error
	GetEventsByGameId(ctx context.Context, gameId string) ([]event.Event, error)
	Insert(ctx context.Context, event event.Event) error
}

type database struct {
	dbClient     *redis.Client
	dbName       string
	dbCollection string
}

func Init(ctx context.Context, cfg config.Config) Database {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.DbURI,
		Password: "",
		DB:       0,
	})

	db := database{
		dbClient: client,
	}

	return &db
}

func (d *database) Disconnect(ctx context.Context) error {
	return d.dbClient.Close()
}

func (d *database) GetEventsByGameId(ctx context.Context, gameId string) ([]event.Event, error) {
	events := make([]event.Event, 0, 10)
	rawEvents, err := d.dbClient.LRange(ctx, gameId, 0, -1).Result()
	if err != nil {
		return []event.Event{}, err
	}

	for _, raw := range rawEvents {
		e, err := event.FromJSON([]byte(raw))
		if err != nil {
			return []event.Event{}, err
		}
		events = append(events, e)
	}

	return events, nil
}

func (d *database) Insert(ctx context.Context, event event.Event) error {
	_, err := d.dbClient.RPush(ctx, event.GameID, event.JsonString()).Result()
	return err
}
