package main

import (
	"github.com/CayenneLow/codenames-eventrouter/config"
	"github.com/CayenneLow/codenames-eventrouter/internal"
)

func main() {
	cfg := config.Init()
	internal.StartServer(cfg)
}
