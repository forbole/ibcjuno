package cmd

import (
	initcmd "github.com/MonikaCat/ibcjuno/cmd/init"
	startcmd "github.com/MonikaCat/ibcjuno/cmd/start/types"
)

// Config represents the general configuration for the commands
type Config struct {
	name        string
	initConfig  *initcmd.Config
	startConfig *startcmd.Config
}

// NewConfig allows to build a new Config instance
func NewConfig(name string) *Config {
	return &Config{
		name: name,
	}
}

// GetName returns the name of the root command
func (c *Config) GetName() string {
	return c.name
}

// GetInitConfig returns the config used during init command
func (c *Config) GetInitConfig() *initcmd.Config {
	if c.initConfig == nil {
		return initcmd.NewConfig()
	}
	return c.initConfig
}

// WithStartConfig sets cfg as the start command configuration
func (c *Config) WithStartConfig(cfg *startcmd.Config) *Config {
	c.startConfig = cfg
	return c
}

// GetStartConfig returns the config used during start command
func (c *Config) GetStartConfig() *startcmd.Config {
	if c.startConfig == nil {
		return startcmd.NewConfig()
	}
	return c.startConfig
}
