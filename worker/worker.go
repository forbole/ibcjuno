package worker

import (
	"fmt"
	"time"

	postgresql "github.com/MonikaCat/ibcjuno/db"

	"github.com/MonikaCat/ibcjuno/utils"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

type Worker struct {
	db postgresql.Database
}

func NewWorker(ctx *Context) Worker {
	return Worker{db: ctx.Database}
}

func (w Worker) StartIBCJuno() {

	if err := w.process(); err != nil {
		go func() {
			log.Info().Msg("error when starting processig token prices")
		}()
	}

}

// process starts the cron job to fetch and store tokens prices every 2 mins
func (w Worker) process() error {
	scheduler := gocron.NewScheduler(time.UTC)

	// Fetch the token prices every 2 mins
	if _, err := scheduler.Every(10).Seconds().Do(func() {
		utils.WatchMethod(w.updatePrice)
	}); err != nil {
		return fmt.Errorf("error while setting up period operations: %s", err)
	}

	scheduler.StartAsync()
	return nil
}

func (w Worker) updatePrice() error {
	log.Info().Msg("updating prices...")
	// Get latest tokens prices
	prices, err := w.db.GetTokenPrices()
	if err != nil {
		return err
	}

	// Save the token prices
	err = w.db.SaveTokensPrices(prices)
	if err != nil {
		return fmt.Errorf("error while saving token prices: %s", err)
	}

	return nil
}

// StoreTokensDetails saves tokens defined inside config.yaml file into database
func (w *Worker) StoreTokensDetails(cfg utils.Config) error {
	for _, coin := range cfg.Tokens.Tokens {
		// Save the coin as a token with its units
		err := w.db.SaveToken(coin)
		if err != nil {
			return fmt.Errorf("error while saving token: %s", err)
		}
	}

	return nil
}
