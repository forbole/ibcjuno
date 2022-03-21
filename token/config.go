package token

import (
	"gopkg.in/yaml.v3"
)

// Config contains the configuration about the tokens
type Config struct {
	Tokens []Token `yaml:"token"`
}

// NewConfig returns a new Config instance
func NewConfig(tokens []Token) *Config {
	return &Config{
		Tokens: tokens,
	}
}

func ParseConfig(bz []byte) (*Config, error) {
	type T struct {
		Config *Config `yaml:"tokens"`
	}
	var cfg T
	err := yaml.Unmarshal(bz, &cfg)
	return cfg.Config, err
}
