package config

import (
	"github.com/CayenneLow/codenames-eventrouter/internal/logger"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	WsHost       string `envconfig:"WS_HOST" default:"localhost"`
	WsPort       string `envconfig:"WS_PORT" default:"8080"`
	LogLevel     string `envconfig:"LOG_LEVEL" default:"INFO"`
	DbURI        string `envconfig:"DB_URI" default:"mongodb://root:example@mongo:27017/"`
	DbName       string `envconfig:"DB_NAME" default:"event_store"`
	DbCollection string `envconfig:"DB_COLLECTION" default:"events"`
}

func Init() Config {
	var cfg Config
	// Initialize ENV
	err := envconfig.Process("codenames-router", &cfg)
	// Initialize logger
	logger.Init(cfg.LogLevel)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Error initializing ENV variables"))
	}

	return cfg
}
