package token_price

import (
	"fmt"
	"strings"

	ibctoken "github.com/forbole/ibcjuno/parser/ibc_token"
	types "github.com/forbole/ibcjuno/types"
)

// GetLatestTokensPrices queries the remote APIs to get the latest prices
func GetLatestTokensPrices(ids []string) ([]types.TokenPrice, error) {
	var tokensPrices []types.MarketTicker
	priceIDs := SplitPriceIDs(ids)

	for _, priceIDSlice := range priceIDs {
		if len(priceIDSlice) == 0 {
			continue
		}

		var prices []types.MarketTicker
		query := fmt.Sprintf("/coins/markets?vs_currency=usd&ids=%s", strings.Join(priceIDSlice, ","))
		err := ibctoken.QueryCoingecko(query, &prices, false)
		if err != nil {
			return nil, err
		}
		tokensPrices = append(tokensPrices, prices...)

	}

	return ConvertCoingeckoPrices(tokensPrices), nil
}

func ConvertCoingeckoPrices(prices []types.MarketTicker) []types.TokenPrice {
	tokenPrices := make([]types.TokenPrice, len(prices))
	for i, price := range prices {
		tokenPrices[i] = types.NewTokenPrice(
			price.CoingeckoID,
			price.Symbol,
			price.Name,
			price.Image,
			price.CurrentPrice,
			price.MarketCap,
			price.MarketCapRank,
			price.FullyDilutedValuation,
			price.TotalVolume,
			price.High24Hrs,
			price.Low24Hrs,
			price.CirculatingSupply,
			price.TotalSupply,
			price.MaxSupply,
			price.ATH,
			price.ATL,
			price.LastUpdated,
		)
	}
	return tokenPrices
}

func SplitPriceIDs(ids []string) [][]string {
	maxBalancesPerSlice := 50
	slices := make([][]string, len(ids)/maxBalancesPerSlice+1)

	sliceIndex := 0
	for index, priceID := range ids {
		slices[sliceIndex] = append(slices[sliceIndex], priceID)

		if index > 0 && index%(maxBalancesPerSlice-1) == 0 {
			sliceIndex++
		}
	}

	return slices
}
