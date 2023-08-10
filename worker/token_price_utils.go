package worker

import (
	"fmt"
	"time"

	tokenprice "github.com/forbole/ibcjuno/parser/token_price"
	"github.com/forbole/ibcjuno/utils"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// StartFetchingPrices starts the cron job that fetches and storees
// tokens prices every 5 minutes
func (w Worker) StartFetchingPrices() error {
	scheduler := gocron.NewScheduler(time.UTC)

	// Fetch the token prices every 5 mins
	if _, err := scheduler.Every(5).Minutes().Do(func() {
		utils.WatchMethod(w.UpdateTokenPrices)
	}); err != nil {
		return fmt.Errorf("error while setting up token prices period operations: %s", err)
	}

	scheduler.StartAsync()

	return nil
}

// UpdateTokenPrices queries latest IBC token prices
// and stores updated values in database
func (w Worker) UpdateTokenPrices() error {
	log.Info().Msg("updating token prices...")

	// Get latest tokens price IDs
	priceIDList, err := w.db.GetTokensPriceID()
	if err != nil {
		return fmt.Errorf("error while getting tokens price ID: %s", err)
	}

	// Get the tokens prices
	prices, err := tokenprice.GetLatestTokensPrices(priceIDList)
	if err != nil {
		return fmt.Errorf("error while getting token prices: %s", err)
	}

	if len(prices) > 0 {
		// Save the token prices
		err = w.db.SaveTokensPrices(prices)
		if err != nil {
			return fmt.Errorf("error while saving token prices: %s", err)
		}
	}

	return nil
}
