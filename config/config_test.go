package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestInit(t *testing.T) {
	cfg := Init()
	f, err := os.Open("./forwarding_rules.yaml")
	var bytes []byte
	f.Read(bytes)
	assert.NoError(t, err)
	var expected map[string](map[string](map[string]([]string)))
	decoder := yaml.NewDecoder(f)
	decoder.Decode(&expected)
	assert.NoError(t, err)
	assert.Equal(t, expected["forwarding_rules"], cfg.ForwardingRules)
}
