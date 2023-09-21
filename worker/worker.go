package worker

import (
	postgresql "github.com/forbole/ibcjuno/database"

	"github.com/rs/zerolog/log"
)

// Worker defines a job consumer that is responsible for getting
// latest IBC token price data and exporting it to a database.
type Worker struct {
	db postgresql.Database
}

// NewWorker allows to create a new Worker implementation.
func NewWorker(ctx *Context) Worker {
	return Worker{db: ctx.Database}
}

// StartWorker starts a worker that triggers periodic operations
// to update the IBC tokens info and price stored in database in given interval.
// It returns an error if any export process fails.
func (w Worker) StartWorker() {

	// Start fetching tokens prices every 5 mins
	if err := w.StartFetchingPrices(); err != nil {
		go func() {
			log.Info().Msg("error when starting processig token prices")
		}()
	}

	// Start fetching IBC token details every day
	if err := w.StartFetchingLatestIBCTokensInfo(); err != nil {
		go func() {
			log.Info().Msg("error when starting processig latest IBC tokens info")
		}()
	}

}
