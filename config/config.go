package config

import (
	"github.com/CayenneLow/codenames-eventrouter/internal/logger"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	WsHost string `envconfig:"WS_HOST" default:"localhost"`
	WsPort string `envconfig:"WS_PORT" default:"8080"`
}

func Init() Config {
	var cfg Config
	// Initialize logger
	logger.Init()
	// Initialize ENV
	err := envconfig.Process("codenames-router", &cfg)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Error initializing ENV variables"))
	}

	return cfg
}
