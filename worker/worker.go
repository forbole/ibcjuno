package worker

import (
	"fmt"
	"time"

	"github.com/MonikaCat/ibc-token/config"
	"github.com/MonikaCat/ibc-token/db"
	"github.com/MonikaCat/ibc-token/utils"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

type (
	Worker struct {
		db *db.Database
	}
)

func NewWorker(db *db.Database) Worker {
	return Worker{db}
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
	if _, err := scheduler.Every(3).Seconds().Do(func() {
		utils.WatchMethod(w.updatePrice)
	}); err != nil {
		return fmt.Errorf("error while setting up period operations: %s", err)
	}

	scheduler.StartAsync()
	return nil
}

func (w Worker) updatePrice() error {

	// Get latest tokens prices
	prices, err := db.GetTokenPrices(w.db)
	if err != nil {
		return err
	}

	// Save the token prices
	err = db.SaveTokensPrices(prices, w.db)
	if err != nil {
		return fmt.Errorf("error while saving token prices: %s", err)
	}

	return nil
}

// StoreTokensDetails saves tokens defined inside config.yaml file into database
func (w *Worker) StoreTokensDetails(cfg config.Config) error {
	for _, coin := range cfg.Tokens.Tokens {
		// Save the coin as a token with its units
		err := db.SaveToken(coin, w.db)
		if err != nil {
			return fmt.Errorf("error while saving token: %s", err)
		}
	}

	return nil
}
