package token_price

import (
	"fmt"
	"math"
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
		err := ibctoken.QueryCoingecko(query, &prices)
		if err != nil {
			return nil, err
		}
		tokensPrices = append(tokensPrices, prices...)
		fmt.Printf("\n\n query %v \n\n", query)

	}

	fmt.Printf("\n\n price ids len %v \n\n", len(ids))

	return ConvertCoingeckoPrices(tokensPrices), nil
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
