package utils

import (
	"os"

	"github.com/forbole/ibcjuno/types"
)

func RemoveDuplicatedTickers(tickers []types.CoinGeckoTicker) []types.CoinGeckoTicker {
	result := []types.CoinGeckoTicker{}

	for i := 0; i < len(tickers); i++ {
		// Check if the ticker has already been added to the result slice
		duplicate := false
		for j := 0; j < len(result); j++ {
			if tickers[i].OriginChainPriceID == result[j].OriginChainPriceID &&
				tickers[i].OriginDenom == result[j].OriginDenom &&
				tickers[i].TargetChainPriceID == result[j].TargetChainPriceID &&
				tickers[i].TargetDenom == result[j].TargetDenom {
				duplicate = true
				break
			}
		}
		// Add the ticker value to the result slice if it's not a duplicate
		if !duplicate {
			result = append(result, tickers[i])
		}
	}
	return result
}

// GetEnvOr returns the value of the ENV variable having the given key, or the provided orValue
func GetEnvOr(envKey string, orValue string) string {
	if envValue := os.Getenv(envKey); envValue != "" {
		return envValue
	}
	return orValue
}
