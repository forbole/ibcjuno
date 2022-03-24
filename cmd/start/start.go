package start

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/MonikaCat/ibcjuno/utils"
	worker "github.com/MonikaCat/ibcjuno/worker"
	"github.com/rs/zerolog/log"

	types "github.com/MonikaCat/ibcjuno/cmd/start/types"
	"github.com/spf13/cobra"
)

var (
	waitGroup sync.WaitGroup
)

// NewStartCmd returns the command that should be run when we want to start parsing a chain state.
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

			return StartParsing(context)
		},
	}
}

// StartParsing represents the function that should be called when the parse command is executed
func StartParsing(ctx *worker.Context) error {

	cfg := utils.Cfg
	workerCount := 1
	workers := make([]worker.Worker, workerCount, workerCount)
	for i := range workers {
		workers[i] = worker.NewWorker(ctx)
	}

	waitGroup.Add(1)

	// Start each blocking worker in a go-routine where the worker consumes jobs
	// off of the export queue.
	for i, w := range workers {
		log.Info().Int("number", i+1).Msg("starting worker...")
		err := w.StoreTokensDetails(cfg)
		if err != nil {
			return err
		}

		go w.StartIBCJuno()
	}

	// listen for and trap any OS signal to gracefully shutdown and exit
	trapSignal(ctx)

	// block main process (signal capture will call WaitGroup's Done)
	waitGroup.Wait()
	return nil
}

// trapSignal will listen for any OS signal and invoke Done on the main
// WaitGroup allowing the main process to gracefully exit.
func trapSignal(ctx *worker.Context) {
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
