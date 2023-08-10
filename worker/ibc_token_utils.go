package worker

import (
	"fmt"
	"time"

	ibctoken "github.com/forbole/ibcjuno/parser/ibc_token"
	"github.com/forbole/ibcjuno/types"
	"github.com/forbole/ibcjuno/utils"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// StartFetchingLatestIBCTokensInfo starts the cron job that fetches and storees
// latest IBC tokens info inside the database once a day
func (w Worker) StartFetchingLatestIBCTokensInfo() error {
	scheduler := gocron.NewScheduler(time.UTC)

	// Setup a cron job to fetch the latest IBC token details every day
	if _, err := scheduler.Every(1).Day().At("00:00").Do(func() {
		utils.WatchMethod(w.QueryAndSaveLatestIBCTokensInfo)
	}); err != nil {
		return fmt.Errorf("error while setting up period operations to fetch latest IBC token details: %s", err)
	}

	scheduler.StartAsync()

	return nil
}

// GetIBCTokensList queries the latest IBC chain list
// and latest IBC token details from the given endpoint
func (w *Worker) GetIBCTokensList() ([]types.IBCTokenUnit, error) {
	// query list of IBC supported networks
	chainList, err := ibctoken.QueryIBCChainList()
	if err != nil {
		log.Error().Err(err).Msg("error while getting IBC chain list")
		return []types.IBCTokenUnit{}, err
	}

	if len(chainList) == 0 {
		panic("IBC chain list is empty")
	}

	// query IBC tokens details for each chain
	tokenList, err := ibctoken.QueryIBCTokensDetails(chainList)
	if err != nil {
		log.Error().Err(err).Msg("error while getting IBC tokens details")
		return nil, err
	}

	return tokenList, nil
}

// QueryAndSaveLatestIBCTokensInfo queries the latest IBC token details
// from the given endpoint and stores them inside the database
func (w *Worker) QueryAndSaveLatestIBCTokensInfo() error {
	log.Info().Msg("getting IBC tokens list...")

	// query the latest IBC tokens list
	ibcTokensList, err := w.GetIBCTokensList()
	if err != nil {
		return fmt.Errorf("error while getting IBC tokens info: %s", err)
	}

	log.Info().Msg("getting IBC tokens details...")

	// store updated IBC tokens list in database
	err = w.db.SaveIBCTokens(ibcTokensList)
	if err != nil {
		return fmt.Errorf("error while saving IBC tokens: %s", err)
	}

	return nil
}
