package start

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/forbole/ibcjuno/utils"
	workerctx "github.com/forbole/ibcjuno/worker"

	"github.com/spf13/cobra"

	types "github.com/forbole/ibcjuno/cmd/start/config"
)

var (
	waitGroup sync.WaitGroup
)

// NewStartCmd returns the command that is run when starting IBCJuno
func NewStartCmd(cmdCfg *types.StartConfig) *cobra.Command {
	return &cobra.Command{
		Use:     "start",
		Short:   "Start IBCJuno price aggregator",
		PreRunE: types.ReadConfigPreRunE(cmdCfg),
		RunE: func(cmd *cobra.Command, args []string) error {
			context, err := types.GetStartContext(utils.Cfg, cmdCfg)
			if err != nil {
				return err
			}

			return IBCJuno(context)
		},
	}
}

// IBCJuno represents the function that is called when
// "start" command is executed
func IBCJuno(ctx *workerctx.Context) error {

	// create worker responsible for fetching latest prices
	worker := workerctx.NewWorker(ctx)
	waitGroup.Add(1)

	// check if refresh ibc tokens been set to true
	// inside config file
	if utils.Cfg.Parser.RefreshIBCTokensOnStart {
		// update IBC tokens details in database to the latest
		// before starting the worker
		err := worker.QueryAndSaveLatestIBCTokensInfo()
		if err != nil {
			return fmt.Errorf("error while saving IBC tokens: %s", err)
		}
	}

	// start worker
	log.Info().Msg("starting worker...")
	go worker.StartWorker()

	// listen for and trap any OS signal to gracefully shutdown and exit
	trapSignal(ctx)

	// block main process (signal capture will call WaitGroup's Done)
	waitGroup.Wait()
	return nil
}

// trapSignal will listen for any OS signal and invoke Close on Database
// and Done on the main WaitGroup allowing the main process to exit gracefully.
func trapSignal(ctx *workerctx.Context) {
	var sigCh = make(chan os.Signal, 1)

	signal.Notify(sigCh, syscall.SIGTERM)
	signal.Notify(sigCh, syscall.SIGINT)

	go func() {
		sig := <-sigCh
		log.Info().Str("signal", sig.String()).Msg("caught signal; IBCJuno is shutting down...")
		defer ctx.Database.Close()
		defer waitGroup.Done()
	}()
}
