package main

import (
	"os"

	"github.com/MonikaCat/ibcjuno/cmd"
	types "github.com/MonikaCat/ibcjuno/cmd/start/types"
)

func main() {
	// IBCJuno config runner
	config := cmd.NewCmdConfig("IBCJuno").StartConfig(types.NewStartConfig())

	// Run the commands and panic if there is any error
	exec := cmd.BuildDefaultCmd(config)
	err := exec.Execute()
	if err != nil {
		os.Exit(1)
	}
}
