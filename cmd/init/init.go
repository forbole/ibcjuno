package init

import (
	"fmt"
	"os"

	config "github.com/forbole/ibcjuno/utils"

	"github.com/spf13/cobra"
)

const (
	replaceFlag = "replace"
)

// NewInitCmd returns the command that should be run to properly initialize IBCJuno config files
func NewInitCmd(cfg *InitConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize configuration files",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create the config path if doesn't exist
			if _, err := os.Stat(config.HomePath); os.IsNotExist(err) {
				err = os.MkdirAll(config.HomePath, os.ModePerm)
				if err != nil {
					return err
				}
			}
			replace, err := cmd.Flags().GetBool(replaceFlag)
			if err != nil {
				return err
			}

			// Get the config file
			configFilePath := config.GetConfigFilePath()
			file, _ := os.Stat(configFilePath)

			// Check if the file exists
			// Replace the file if --replace flag is used
			if file != nil && !replace {
				return fmt.Errorf(
					"configuration file already exists at %s. If you wish to overwrite it, use the --%s flag",
					configFilePath, replaceFlag)
			}
			// Get the config from the flags
			yamlCfg := cfg.GetConfigCreator()(cmd)

			return config.Write(yamlCfg, configFilePath)
		},
	}

	cmd.Flags().Bool(replaceFlag, false, "overrides any existing configuration file")

	return cmd
}
