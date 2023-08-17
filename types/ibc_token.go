package types

import "time"

type ChainRegistryList struct {
	Name string `json:"name" yaml:"name"`
}

type ChainRegistryAssetsList struct {
	Assets []ChainRegistryAsset `json:"assets"`
}

type ChainRegistryAsset struct {
	DenomUnits  []DenomUnit `json:"denom_units" yaml:"denom_units"`
	Base        string      `json:"base" yaml:"base"`
	Name        string      `json:"name" yaml:"name"`
	Display     string      `json:"display" yaml:"display"`
	Symbol      string      `json:"symbol" yaml:"symbol"`
	CoingeckoID string      `json:"coingecko_id" yaml:"coingecko_id"`
}

type DenomUnit struct {
	Denom    string `json:"denom" yaml:"denom"`
	Exponent int    `json:"exponent" yaml:"exponent"`
}

type CoinGeckoTokenDetailsResponse struct {
	Name    string            `json:"name" yaml:"name"`
	Tickers []CoinGeckoTicker `json:"tickers" yaml:"tickers"`
}

type CoinGeckoTicker struct {
	Denom       string    `json:"base" yaml:"base"`
	OriginChain string    `json:"coin_id" yaml:"coin_id"`
	TargetDenom string    `json:"target" yaml:"target"`
	TargetChain string    `json:"target_coin_id" yaml:"target_coin_id"`
	IsStale     bool      `json:"is_stale" yaml:"is_stale"`
	TradeURL    string    `json:"trade_url" yaml:"trade_url"`
	Timestamp   time.Time `json:"timestamp" yaml:"timestamp"`
}

type IBCToken struct {
	DenomUnits  []DenomUnit       `json:"denom_units" yaml:"denom_units"`
	Base        string            `json:"base" yaml:"base"`
	Name        string            `json:"name" yaml:"name"`
	Display     string            `json:"display" yaml:"display"`
	Symbol      string            `json:"symbol" yaml:"symbol"`
	CoingeckoID string            `json:"coingecko_id" yaml:"coingecko_id"`
	Tickers     []CoinGeckoTicker `json:"tickers" yaml:"tickers"`
}

func NewIBCToken(denomUnits []DenomUnit, base, name, display, symbol,
	coingeckoID string, tickers []CoinGeckoTicker) IBCToken {
	return IBCToken{
		DenomUnits:  denomUnits,
		Base:        base,
		Name:        name,
		Display:     display,
		Symbol:      symbol,
		CoingeckoID: coingeckoID,
		Tickers:     tickers,
	}
}
