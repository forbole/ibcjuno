package start

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/MonikaCat/ibcjuno/utils"
	workerctx "github.com/MonikaCat/ibcjuno/worker"
	"github.com/rs/zerolog/log"

	types "github.com/MonikaCat/ibcjuno/cmd/start/types"
	"github.com/spf13/cobra"
)

var (
	waitGroup sync.WaitGroup
)

// NewStartCmd returns the command that should be run when starting IBCJuno
func NewStartCmd(cmdCfg *types.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "start",
		Short:   "Start IBCJuno price aggregator",
		PreRunE: types.ReadConfigPreRunE(cmdCfg),
		RunE: func(cmd *cobra.Command, args []string) error {
			context, err := types.GetStartContext(utils.Cfg, cmdCfg)
			if err != nil {
				return err
			}

			return StartFetchingPrices(context)
		},
	}
}

// StartFetchingPrices represents the function that is called when
// start command is executed
func StartFetchingPrices(ctx *workerctx.Context) error {
	log.Info().Msg("starting worker...")

	// Create worker responsible for fetching latest prices
	worker := workerctx.NewWorker((ctx))
	waitGroup.Add(1)

	// Get the config
	cfg := utils.Cfg

	// Store tokens defined in config file inside the database
	err := worker.StoreTokensDetails(cfg)
	if err != nil {
		return err
	}

	// Start IBCJuno
	go worker.StartIBCJuno()

	// listen for and trap any OS signal to gracefully shutdown and exit
	trapSignal(ctx)

	// block main process (signal capture will call WaitGroup's Done)
	waitGroup.Wait()
	return nil
}

// trapSignal will listen for any OS signal and invoke Close on Database
// and Done on the main WaitGroup allowing the main process to gracefully exit.
func trapSignal(ctx *workerctx.Context) {
	var sigCh = make(chan os.Signal)

	signal.Notify(sigCh, syscall.SIGTERM)
	signal.Notify(sigCh, syscall.SIGINT)

	go func() {
		sig := <-sigCh
		log.Info().Str("signal", sig.String()).Msg("caught signal; shutting down...")
		defer ctx.Database.Close()
		defer waitGroup.Done()
	}()
}
