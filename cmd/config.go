package cmd

import (
	initcmd "github.com/forbole/ibcjuno/cmd/init"
	startcmd "github.com/forbole/ibcjuno/cmd/start/config"
)

// CmdConfig represents configuration for "init" and "start" commands
type CmdConfig struct {
	name        string
	initConfig  *initcmd.InitConfig
	startConfig *startcmd.StartConfig
}

// NewCmdConfig allows to build a new CmdConfig instance
func NewCmdConfig(name string) *CmdConfig {
	return &CmdConfig{
		name: name,
	}
}

// GetName returns the name of the root command
func (c *CmdConfig) GetName() string {
	return c.name
}

// GetInitConfig returns the config used during init command
func (c *CmdConfig) GetInitConfig() *initcmd.InitConfig {
	if c.initConfig == nil {
		return initcmd.NewInitConfig()
	}
	return c.initConfig
}

// StartConfig sets cfg as the start command configuration
func (c *CmdConfig) StartConfig(cfg *startcmd.StartConfig) *CmdConfig {
	c.startConfig = cfg
	return c
}

// GetStartConfig returns the config used during start command
func (c *CmdConfig) GetStartConfig() *startcmd.StartConfig {
	if c.startConfig == nil {
		return startcmd.NewStartConfig()
	}
	return c.startConfig
}
