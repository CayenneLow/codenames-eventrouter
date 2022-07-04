package config

import (
	"log"
	"os"

	"github.com/CayenneLow/codenames-eventrouter/internal/logger"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ForwardingRules map[string](map[string]([]string)) `yaml:"forwarding_rules"`
	WsHost          string                             `envconfig:"WS_HOST" default:"localhost"`
	WsPort          string                             `envconfig:"WS_PORT" default:"8080"`
}

func Init() Config {
	// Initialize Forwarding Rules
	f, err := os.Open("config/forwarding_rules.yaml")
	if err != nil {
		log.Fatal(errors.Wrap(err, "Error opening forwarding rules file"))
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Error initializing forwarding rules"))
	}

	// Initialize envconfig
	err = envconfig.Process("codenames-router", &cfg)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Error initializing ENV variables"))
	}

	// Initialize logger
	logger.Init()

	return cfg
}

func (c *Config) GetReceivers(eventType string) []string {
	return c.ForwardingRules[eventType]["receivers"]
}

func (c *Config) GetAcknowledgers(eventType string) []string {
	return c.ForwardingRules[eventType]["acknowledgers"]
}
