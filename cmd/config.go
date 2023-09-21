package cmd

import (
	initcmd "github.com/forbole/ibcjuno/cmd/init"
	startcmd "github.com/forbole/ibcjuno/cmd/start/config"
)

// Config represents configuration for "init" and "start" commands
type Config struct {
	name        string
	initConfig  *initcmd.Config
	startConfig *startcmd.StartConfig
}

// NewCmdConfig allows to build a new Config instance
func NewCmdConfig(name string) *Config {
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
		return initcmd.NewInitConfig()
	}
	return c.initConfig
}

// StartConfig sets cfg as the start command configuration
func (c *Config) StartConfig(cfg *startcmd.StartConfig) *Config {
	c.startConfig = cfg
	return c
}

// GetStartConfig returns the config used during start command
func (c *Config) GetStartConfig() *startcmd.StartConfig {
	if c.startConfig == nil {
		return startcmd.NewStartConfig()
	}
	return c.startConfig
}
