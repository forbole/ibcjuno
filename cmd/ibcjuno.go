package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/MonikaCat/ibc-token/config"
	"github.com/MonikaCat/ibc-token/db"
	worker "github.com/MonikaCat/ibc-token/worker"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	logLevelJSON = "json"
	logLevelText = "text"
)

var (
	logLevel  string
	logFormat string
	wg        sync.WaitGroup
)

var rootCmd = &cobra.Command{
	Use:   "ibcjuno [config-file]",
	Args:  cobra.ExactArgs(1),
	Short: "IBCJuno is a IBC tokens price aggregator and exporter",
	Long: `IBCJuno IBC tokens price aggregator. It queries latest IBC tokens prices and indexes it using PostgreSQL database. IBCJuno is meat to run with a GraphQL layer on top 
	to ease the ability for developers and downstream clients to query the latest price of any IBC token.`,
	RunE: ibcJunoCmdHandler,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", zerolog.InfoLevel.String(), "logging level")
	rootCmd.PersistentFlags().StringVar(&logFormat, "log-format", logLevelJSON, "logging format; must be either json or text")

	rootCmd.AddCommand(getVersionCmd())
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ibcJunoCmdHandler(cmd *cobra.Command, args []string) error {
	logLvl, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return err
	}

	zerolog.SetGlobalLevel(logLvl)

	switch logFormat {
	case logLevelJSON:
		// JSON is the default logging format

	case logLevelText:
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	default:
		return fmt.Errorf("invalid logging format: %s", logFormat)
	}

	cfgFile := args[0]
	cfg := config.ParseConfig(cfgFile)
	workerCount := 1

	db, err := db.OpenDB(cfg)
	if err != nil {
		return errors.Wrap(err, "failed to open database connection")
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		return errors.Wrap(err, "failed to ping database")
	}

	workers := make([]worker.Worker, workerCount)
	for i := range workers {
		workers[i] = worker.NewWorker(db)
	}

	wg.Add(1)

	// Start each blocking worker in a go-routine where the worker consumes jobs
	// off of the export queue.
	for i, w := range workers {
		log.Info().Int("number", i+1).Msg("starting worker...")
		err = w.StoreTokensDetails(cfg)
		if err != nil {
			return err
		}

		go w.StartIBCJuno()
	}

	// listen for and trap any OS signal to gracefully shutdown and exit
	trapSignal()

	// block main process (signal capture will call WaitGroup's Done)
	wg.Wait()
	return nil
}

// trapSignal will listen for any OS signal and invoke Done on the main
// WaitGroup allowing the main process to gracefully exit.
func trapSignal() {
	var sigCh = make(chan os.Signal)

	signal.Notify(sigCh, syscall.SIGTERM)
	signal.Notify(sigCh, syscall.SIGINT)

	go func() {
		sig := <-sigCh
		log.Info().Str("signal", sig.String()).Msg("caught signal; shutting down...")
		defer wg.Done()
	}()
}
