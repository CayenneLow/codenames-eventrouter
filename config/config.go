package config

import (
	"log"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ForwardingRules map[string](map[string]([]string)) `yaml:"forwarding_rules"`
}

func Init() Config {
	// Initialize Forwarding Rules
	f, err := os.Open("forwarding_rules.yaml")
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

	return cfg
}

func (c *Config) GetReceivers(eventType string) []string {
	return c.ForwardingRules[eventType]["receivers"]
}

func (c *Config) GetAcknowledgers(eventType string) []string {
	return c.ForwardingRules[eventType]["acknowledgers"]
}
