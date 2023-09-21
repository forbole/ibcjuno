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
	if len(utils.Cfg.API.ChainRegistryAPIURL) == 0 {
		panic("Chain registry url inside config.yaml file is empty")
	}

	resp, err := http.Get(fmt.Sprintf("%s/contents", utils.Cfg.API.ChainRegistryAPIURL))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error: %s ", resp.Status)
	}

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

	return chainList, nil
}

// QueryIBCAssetsDetailsFromChainRegistry queries IBC token details for each chain
// from chain registry
func QueryIBCAssetsDetailsFromChainRegistry(chainList []string) ([]types.ChainRegistryAsset, error) {
	var chainRegistryAsset []types.ChainRegistryAsset
	var tokenList []types.ChainRegistryAsset

	// panic if chain registry url is empty
	if len(utils.Cfg.API.ChainRegistryRawURL) == 0 {
		panic("Chain registry url inside config.yaml file is empty")
	}

	// query each chain IBC token details
	for _, network := range chainList {
		var ibcDenom types.ChainRegistryAssetsList
		url := fmt.Sprintf("%s/master/%s/assetlist.json", utils.Cfg.API.ChainRegistryRawURL, network)
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Info().Msgf("assets.json file not found for %s chain, skipping...", network)
			continue
		}

		bz, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error while reading %s chain response body: %s", network, err)
		}

		err = json.Unmarshal(bz, &ibcDenom)
		if err != nil {
			return nil, fmt.Errorf("error while unmarshaling %s chain response body: %s",network, err)
		}

		chainRegistryAsset = append(chainRegistryAsset, ibcDenom.Assets...)
	}

	for _, asset := range chainRegistryAsset {
		if len(asset.CoingeckoID) > 0 {
			tokenList = append(tokenList, asset)
		}
	}

	return tokenList, nil
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
		time.Sleep(6 * time.Second)
		return nil
	}
	if resp.StatusCode == 429 {
		log.Error().Msg("error 429: too many requests... will try to refetch again...")
		time.Sleep(20 * time.Second)
		return fmt.Errorf("error 429")
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

	// wait for 6 sec to avoid 429 error
	time.Sleep(6 * time.Second)

	return nil
}
