package config

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/forbole/ibcjuno/utils"
)

// ReadConfigPreRunE represents a Cobra cmd function allowing to read the config before executing the command itself
func ReadConfigPreRunE(cfg *StartConfig) utils.CobraCmdFunc {
	return func(_ *cobra.Command, _ []string) error {
		return UpdateGlobalConfig(cfg)
	}
}

// ReadConfig allows to read the configuration using the provided config
func ReadConfig(cfg *StartConfig) (utils.Config, error) {
	file := utils.GetConfigFilePath()

	// Ensure the path and config file exists
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return utils.Config{}, fmt.Errorf("config file does not exist. Make sure you have run the init command")
	}

	// Read the config file
	return utils.Read(file, cfg.GetConfigParser())
}

// UpdateGlobalConfig parses the configuration file
// and sets it as global configuration
func UpdateGlobalConfig(cfg *StartConfig) error {
	junoCfg, err := ReadConfig(cfg)
	if err != nil {
		return err
	}

	// Set the global configuration
	utils.Cfg = junoCfg
	return nil
}
