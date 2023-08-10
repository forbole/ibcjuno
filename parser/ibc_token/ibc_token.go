package ibc_tokens

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	types "github.com/forbole/ibcjuno/types"
	"github.com/forbole/ibcjuno/utils"
	"github.com/rs/zerolog/log"
)

// QueryIBCChainList queries the list of IBC supported chains
func QueryIBCChainList() ([]string, error) {
	var ptr []string

	resp, err := http.Get(utils.Cfg.API.SupportedChainsURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("error while reading response body: ")
		return nil, err
	}

	err = json.Unmarshal(bz, &ptr)
	if err != nil {
		log.Error().Err(err).Msg("error while unmarshaling response body: ")
		return nil, err
	}

	return ptr, nil
}

// QueryIBCTokensDetails queries IBC token details for each chain
func QueryIBCTokensDetails(chainList []string) ([]types.IBCTokenUnit, error) {
	var tokenList []types.IBCTokenUnit

	// query each chain IBC token details
	for _, network := range chainList {
		var ibcDenom []types.IBCTokenUnit
		url := fmt.Sprintf(utils.Cfg.API.ChainAssetsURL, network)
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

		tokenList = append(tokenList, ibcDenom...)
	}

	return tokenList, nil
}
