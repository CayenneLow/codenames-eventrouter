package main

import (
	"context"
	"time"

	"github.com/CayenneLow/codenames-eventrouter/config"
	"github.com/CayenneLow/codenames-eventrouter/internal"
	"github.com/CayenneLow/codenames-eventrouter/internal/database"
)

func main() {
	cfg := config.Init()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db := database.Init(ctx, cfg)
	defer db.Disconnect(ctx)
	internal.StartServer(cfg, db)
}
