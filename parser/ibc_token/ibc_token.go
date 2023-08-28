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
			asset.Name = ParseChainName(asset.Name)
			tokenList = append(tokenList, asset)
		}
	}

	return tokenList, nil
}

// QueryCoingecko queries the CoinGecko APIs for the given endpoint
func QueryCoingecko(endpoint string, ptr interface{}, queryIBCToken bool) error {
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
	}
	if resp.StatusCode == 429 {
		log.Error().Msg("error 429: too many requests... will try to refetch again...")
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

	// wait for 15 seconds if querying IBC token details
	// to minimise the 429 error chances
	if queryIBCToken {
		// wait for 15 seconds
		time.Sleep(15 * time.Second)
	}

	return nil
}

// Parse chain name to the one returned by coingecko
// to enable tokens relationship
func ParseChainName(chainName string) string {

	switch chainName {
	case "Cosmos Hub Atom":
		return "Cosmos Hub"
	case "Canto":
		return "CANTO"
	case "fetch-ai":
		return "Fetch.ai"
	case "FirmaChain":
		return "Firmachain"
	case "JunoSwap":
		return "JUNO"
	case "AIOZ":
		return "AIOZ Network"
	case "cheqd":
		return "CHEQD Network"
	case "Carbon":
		return "Carbon Protocol"
	case "Carbon USD Coin":
		return "Carbon USD"
	case "Aura":
		return "Aura Network"
	case "Crescent":
		return "Crescent Network"
	case "Comdex":
		return "COMDEX"
	case "Chihuahua":
		return "Chihuahua Chain"
	case "Jackal":
		return "Jackal Protocol"
	case "CMST":
		return "Composite"
	case "Arable USD":
		return "Arable Protocol"
	case "Neta":
		return "NETA"
	case "LORE":
		return "Gitopia"
	case "USD Coin":
		return "Bridged USD Coin (Axelar)"
	case "Bonded Crescent":
		return "Liquid Staking Crescent"
	case "Marble":
		return "Marble Dao"
	case "OKExChain":
		return "OKT Chain"
	case "Kuji":
		return "Kujira"
	case "NYM":
		return "Nym"
	case "MediBloc":
		return "Medibloc"
	case "Hard":
		return "Kava Lend"
	case "Ki":
		return "KI"
	case "Regen Network":
		return "Regen"
	case "Nom":
		return "Onomy Protocol"
	case "Mars":
		return "Mars Protocol"
	case "MNTA":
		return "MantaDAO"
	case "Whale":
		return "White Whale"
	case "Swap":
		return "Kava Swap"
	case "Rizon Chain":
		return "RIZON"
	case "ODIN":
		return "Odin Protocol"
	case "Realio Network":
		return "Realio"
	case "PSTAKE staked ATOM":
		return "stkATOM"
	case "Nature Carbon Ton":
		return "Toucan Protocol: Nature Carbon Tonne"
	case "DARC":
		return "Konstellation"
	case "Cacao":
		return "Maya Protocol"
	case "MEME":
		return "Meme Network"
	case "Lum":
		return "Lum Network"
	case "Flix":
		return "OmniFlix Network"
	case "Loop Finance":
		return "LOOP"
	case "Luna Classic":
		return "Terra Luna Classic"
	case "Secret Network":
		return "Secret"
	case "Somm":
		return "Sommelier"
	case "FIS":
		return "Stafi"
	case "Sifchain Rowan":
		return "Sifchain"
	case "stATOM":
		return "Stride Staked Atom"
	case "Unification Network":
		return "Unification"
	case "Xpla":
		return "XPLA"
	case "ERIS Amplified LUNA":
		return "Eris Amplified Luna"
	default:
		return chainName
	}

}
