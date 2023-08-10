package token_price

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"

	types "github.com/forbole/ibcjuno/types"
	"github.com/forbole/ibcjuno/utils"
	"github.com/rs/zerolog/log"
)

// GetLatestTokensPrices queries the remote APIs to get the latest prices
func GetLatestTokensPrices(ids []string) ([]types.TokenPrice, error) {
	var prices []types.MarketTicker
	query := fmt.Sprintf("/coins/markets?vs_currency=usd&ids=%s", strings.Join(ids, ","))
	err := queryCoinGecko(query, &prices)
	if err != nil {
		return nil, err
	}

	return ConvertCoingeckoPrices(prices), nil
}

func ConvertCoingeckoPrices(prices []types.MarketTicker) []types.TokenPrice {
	tokenPrices := make([]types.TokenPrice, len(prices))
	for i, price := range prices {
		tokenPrices[i] = types.NewTokenPrice(
			price.Symbol,
			price.CurrentPrice,
			int64(math.Trunc(price.MarketCap)),
			price.LastUpdated,
		)
	}
	return tokenPrices
}

// queryCoinGecko queries the CoinGecko APIs for the given endpoint
func queryCoinGecko(endpoint string, ptr interface{}) error {
	resp, err := http.Get(utils.Cfg.API.CoingeckoURL + endpoint)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	bz, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("error while reading response body: ")
		return err
	}

	err = json.Unmarshal(bz, &ptr)
	if err != nil {
		log.Error().Err(err).Msg("error while unmarshaling response body: ")
		return err
	}

	return nil
}
