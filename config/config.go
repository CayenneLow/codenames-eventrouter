package config

import (
	"os"

	"github.com/CayenneLow/codenames-eventrouter/internal/logger"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ForwardingRules     map[string](map[string]([]string)) `yaml:"forwarding_rules"`
	ForwardingRulesPath string                             `envconfig:"FORWARDING_RULES_PATH"`
	WsHost              string                             `envconfig:"WS_HOST" default:"localhost"`
	WsPort              string                             `envconfig:"WS_PORT" default:"8080"`
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

	// Initialize Forwarding Rules
	log.Debugf("Forwarding Rules Path: %s", cfg.ForwardingRulesPath)
	log.Debugf("Config: %v", cfg)
	f, err := os.Open(cfg.ForwardingRulesPath)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Error opening forwarding rules file"))
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Error initializing forwarding rules"))
	}

	return cfg
}

func (c *Config) GetReceivers(eventType string) []string {
	return c.ForwardingRules[eventType]["receivers"]
}

func (c *Config) GetAcknowledgers(eventType string) []string {
	return c.ForwardingRules[eventType]["acknowledgers"]
}
