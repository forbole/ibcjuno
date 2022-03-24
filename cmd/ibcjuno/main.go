package main

import (
	"os"

	"github.com/MonikaCat/ibcjuno/cmd"
	types "github.com/MonikaCat/ibcjuno/cmd/start/types"
)

func main() {
	// IBCJuno config runner
	config := cmd.NewConfig("ibcjuno").WithStartConfig(types.NewConfig())

	// Run the commands and panic if there is any error
	exec := cmd.BuildDefaultCmd(config)
	err := exec.Execute()
	if err != nil {
		os.Exit(1)
	}
}
