package cmd

import (
	"fmt"
	"os"
	"path"

	initcmd "github.com/MonikaCat/ibcjuno/cmd/init"
	startcmd "github.com/MonikaCat/ibcjuno/cmd/start"

	utils "github.com/MonikaCat/ibcjuno/utils"

	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
)

var (
	FlagHome = "home"
)

// BuildDefaultCmd allows to build an Executor containing a root command that
// has the name, description, default version, init and start commands implementations
func BuildDefaultCmd(config *Config) cli.Executor {
	rootCmd := RootCmd(config.GetName())

	rootCmd.AddCommand(
		getVersionCmd(),
		initcmd.NewInitCmd(config.GetInitConfig()),
		startcmd.NewStartCmd(config.GetStartConfig()),
	)

	return PrepareRootCmd(config.GetName(), rootCmd)
}

// RootCmd allows to build the default root command having the given name
func RootCmd(name string) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("%s is a IBC price aggregator and exporter", name),
		Long: fmt.Sprintf(`%s is a IBC price aggregator and exporter. It queries the latest IBC tokens prices 
and stores it inside PostgreSQL database. %s is meant to run with a GraphQL layer on top to ease the ability for 
developers and downstream clients to query the latest price of any IBC token.`, name, name),
	}
}

// PrepareRootCmd prepares the given command binding all the viper flags
func PrepareRootCmd(name string, cmd *cobra.Command) cli.Executor {
	cmd.PersistentPreRunE = utils.ConcatCobraCmdFuncs(
		utils.BindFlagsLoadViper,
		setupHome,
		cmd.PersistentPreRunE,
	)

	home, _ := os.UserHomeDir()
	defaultConfigPath := path.Join(home, fmt.Sprintf(".%s", name))
	cmd.PersistentFlags().String(FlagHome, defaultConfigPath, "Set the home folder of the application, where all files will be stored")

	return cli.Executor{Command: cmd, Exit: os.Exit}
}

// setupHome sets up home directory of the root command
func setupHome(cmd *cobra.Command, _ []string) error {
	utils.HomePath, _ = cmd.Flags().GetString(FlagHome)
	return nil
}
