package ibc_tokens

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	types "github.com/forbole/ibcjuno/types"
	"github.com/forbole/ibcjuno/utils"
	"github.com/rs/zerolog/log"
)

// QueryIBCChainListFromChainRegistry queries the list of IBC supported chains
func QueryIBCChainListFromChainRegistry() ([]string, error) {
	var chainList []string
	var supportedChains []types.ChainRegistryList

	// panic if chain registry url is empty
	if len(utils.Cfg.API.ChainRegistryURL) == 0 {
		panic("Chain registry url inside config.yaml file is empty")
	}

	resp, err := http.Get(utils.Cfg.API.ChainRegistryURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("error while reading response body: ")
		return nil, err
	}

	err = json.Unmarshal(bz, &supportedChains)
	if err != nil {
		log.Error().Err(err).Msg("error while unmarshaling response body: ")
		return nil, err
	}

	for _, i := range supportedChains {
		if strings.Contains(i.Name, ".") || strings.Contains(i.Name, "_") ||
			strings.Contains(i.Name, "LICENSE") || strings.Contains(i.Name, "testnets") {
			// skip useless values
			continue
		}

		chainList = append(chainList, i.Name)
	}

	fmt.Printf("\n chainList %v \n", chainList)
	fmt.Printf("\n chainList length %v \n", len(chainList))

	return chainList, nil
}

// QueryIBCAssetsDetailsFromChainRegistry queries IBC token details for each chain
// from chain registry
func QueryIBCAssetsDetailsFromChainRegistry(chainList []string) ([]types.ChainRegistryAsset, error) {
	var chainRegistryAsset []types.ChainRegistryAsset
	var tokenList []types.ChainRegistryAsset

	// panic if chain registry assets url is empty
	if len(utils.Cfg.API.ChainRegistryAssetsURL) == 0 {
		panic("Chain registry assets url inside config.yaml file is empty")
	}

	// query each chain IBC token details
	for _, network := range chainList {
		var ibcDenom types.ChainRegistryAssetsList
		url := fmt.Sprintf(utils.Cfg.API.ChainRegistryAssetsURL, network)
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		bz, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error().Err(err).Msg("error while reading response body: ")
		}

		err = json.Unmarshal(bz, &ibcDenom)
		if err != nil {
			log.Info().Msgf("assets.json file not found for %s chain, skipping...", network)
		}

		chainRegistryAsset = append(chainRegistryAsset, ibcDenom.Assets...)
	}

	for _, asset := range chainRegistryAsset {
		if len(asset.CoingeckoID) > 0 {
			tokenList = append(tokenList, asset)
		}
	}

	// priceIDList := utils.RemoveDuplicates(tokenList)

	return tokenList, nil
}

// QueryCoinGeckoForIBCTokensDetails queries the remote APIs to get the latest IBC tokens details
func QueryCoinGeckoForIBCTokensDetails(ids []types.ChainRegistryAsset) ([]types.IBCToken, error) {
	var ibcToken []types.IBCToken
	var missedCoingeckoTokens []types.ChainRegistryAsset

	for i, index := range ids {
		log.Info().Msgf("processing %s network... %d/%d ", index.Name, i+1, len(ids))

		if len(index.CoingeckoID) == 0 {
			continue
		}

		var tokenDetails types.CoinGeckoTokenDetailsResponse
		query := fmt.Sprintf("/coins/%s/tickers", index.CoingeckoID)
		err := QueryCoingecko(query, &tokenDetails)
		if err != nil {
			time.Sleep(5 * time.Second)
			missedCoingeckoTokens = append(missedCoingeckoTokens, index)
		}
		if len(tokenDetails.Tickers) > 0 {
			ibcToken = append(ibcToken, types.NewIBCToken(index.DenomUnits, index.Base, index.Name, index.Display, index.Symbol, index.CoingeckoID, tokenDetails.Tickers))
		}
	}

	fmt.Printf("\n\n missedCoingeckoTokens %v \n\n", missedCoingeckoTokens)

	if len(missedCoingeckoTokens) > 0 {
		log.Info().Msg("*** Refetching previously skipped tokens due to 429 error... ***")
		_, err := QueryCoinGeckoForIBCTokensDetails(missedCoingeckoTokens)
		if err != nil {
			return nil, err
		}

	} else {
		log.Info().Msg("*** Finished processing all networks... Success! ***")

	}

	return ibcToken, nil
}

// QueryCoingecko queries the CoinGecko APIs for the given endpoint
func QueryCoingecko(endpoint string, ptr interface{}) error {
	// panic if coingecko url is empty
	if len(utils.Cfg.API.CoingeckoURL) == 0 {
		panic("Coingecko url inside config.yaml file is empty")
	}

	resp, err := http.Get(utils.Cfg.API.CoingeckoURL + endpoint)
	if err != nil {
		return fmt.Errorf("error while querying coingecko endpoint: %s ", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		log.Info().Msgf("error 404: %s info not found... skipping...", endpoint)
		return nil
	} else if resp.StatusCode == 429 {
		return fmt.Errorf("error 429: too many requests... will try to refetch again... ")
	}

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error while reading response body: %s ", err)
	}

	if len(bz) == 0 {
		return fmt.Errorf("error: response body is 0 bytes")
	}

	err = json.Unmarshal(bz, &ptr)
	if err != nil {
		return fmt.Errorf("error while unmarshaling response body: %s", err)

	}

	// wait for 1 second
	time.Sleep(1 * time.Second)

	return nil
}
