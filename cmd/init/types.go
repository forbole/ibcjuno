package init

import (
	"github.com/spf13/cobra"

	config "github.com/forbole/ibcjuno/utils"
)

// ConfigCreator represents a function that builds a Config instance from the flags that have been specified by the
// user inside the given command.
type ConfigCreator = func(cmd *cobra.Command) config.Config

// DefaultConfigCreator represents the default configuration creator that builds a Config instance using the values
// specified using the default flags.
func DefaultConfigCreator(_ *cobra.Command) config.Config {
	return config.DefaultConfig()
}

// InitConfig contains the configuration data for "init" command
type InitConfig struct {
	createConfig ConfigCreator
}

// NewInitConfig allows to build a new InitConfig instance
func NewInitConfig() *InitConfig {
	return &InitConfig{}
}

// GetConfigCreator return the function that should be run to create a configuration from a set of
// flags specified by the user with the "init" command
func (c *InitConfig) GetConfigCreator() ConfigCreator {
	if c.createConfig == nil {
		return DefaultConfigCreator
	}
	return c.createConfig
}
