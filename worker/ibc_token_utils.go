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
func (w *Worker) GetIBCTokensList() ([]types.ChainRegistryAsset, error) {
	log.Info().Msg("*** Getting IBC chain list from chain registry... ***")

	// query list of IBC supported networks
	chainList, err := ibctoken.QueryIBCChainListFromChainRegistry()
	if err != nil {
		log.Error().Err(err).Msg("error while getting IBC chain list from chain registry")
		return []types.ChainRegistryAsset{}, err
	}

	// panic if queried IBC chain list from chain registry is empty
	if len(chainList) == 0 {
		panic("IBC chain list is empty")
	}

	log.Info().Msg("*** Getting IBC tokens assets info from chain registry... ***")

	// query IBC tokens assets info for each chain
	ibcAssetsDetails, err := ibctoken.QueryIBCAssetsDetailsFromChainRegistry(chainList)
	if err != nil {
		log.Error().Err(err).Msg("error while getting IBC assets details from chain registry")
		return nil, err
	}

	return ibcAssetsDetails, nil
}

// QueryAndSaveLatestIBCTokensInfo queries the latest IBC token details
// from the given endpoint and stores them inside the database
func (w *Worker) QueryAndSaveLatestIBCTokensInfo() error { // start
	log.Info().Msg("*** Getting IBC tokens list... ***")

	// query the latest IBC tokens list
	ibcTokenAssets, err := w.GetIBCTokensList() // above this ok
	if err != nil {
		return fmt.Errorf("error while getting IBC tokens list: %s", err)
	}

	// store updated tokens list in database
	err = w.db.SaveTokens(ibcTokenAssets)
	if err != nil {
		return fmt.Errorf("error while saving IBC tokens in db: %s", err)
	}

	log.Info().Msg("*** Getting IBC tokens details... ***")

	// query the latest IBC tokens details
	err = w.QueryCoinGeckoForIBCTokensDetails(ibcTokenAssets)
	if err != nil {
		return fmt.Errorf("error while getting IBC tokens info: %s", err)
	}

	return nil
}

// QueryCoinGeckoForIBCTokensDetails queries the remote APIs to get the latest IBC tokens details
// and stores updated values in database
func (w *Worker) QueryCoinGeckoForIBCTokensDetails(ids []types.ChainRegistryAsset) error {
	var missedCoingeckoTokens []types.ChainRegistryAsset
	tickerCount := 0
	for i, index := range ids {
		log.Info().Msgf("processing %s network... %d/%d ", index.Name, i+1, len(ids))

		if len(index.CoingeckoID) == 0 {
			continue
		}

		var tokenDetails types.CoinGeckoTokenDetailsResponse
		query := fmt.Sprintf("/coins/%s/tickers", index.CoingeckoID)
		err := ibctoken.QueryCoingecko(query, &tokenDetails, true)
		if err != nil {
			// time.Sleep(10 * time.Second)
			missedCoingeckoTokens = append(missedCoingeckoTokens, index)
		}
		if len(tokenDetails.Tickers) > 0 {
			// store updated IBC tokens list in database
			err = w.db.SaveIBCToken(types.NewIBCToken(index.DenomUnits, index.Base, index.Name, index.Display, index.Symbol, index.CoingeckoID, tokenDetails.Tickers))
			if err != nil {
				return fmt.Errorf("error while saving IBC tokens in db: %s", err)
			}
		}
	}

	if len(missedCoingeckoTokens) > 0 {
		// time.Sleep(30 * time.Second) // disable for testing
		log.Info().Msgf("*** Refetching previously skipped tokens due to 429 error... total %v ***", len(missedCoingeckoTokens))
		err := w.QueryCoinGeckoForIBCTokensDetails(missedCoingeckoTokens)
		if err != nil {
			return err
		}
	} else {
		log.Info().Msg("*** Finished processing all networks... Success! ***")

	}

	fmt.Printf("\n\n tickerCount %d \n\n", tickerCount)

	return nil
}
