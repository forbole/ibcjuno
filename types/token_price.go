package types

import "time"

// MarketTicker contains the current market data for a single token
type MarketTicker struct {
	CoingeckoID           string    `json:"id"`
	Symbol                string    `json:"symbol"`
	Name                  string    `json:"name"`
	Image                 string    `json:"image"`
	CurrentPrice          float64   `json:"current_price"`
	MarketCap             float64   `json:"market_cap"`
	MarketCapRank         int64     `json:"market_cap_rank"`
	FullyDilutedValuation float64   `json:"fully_diluted_valuation"`
	TotalVolume           float64   `json:"total_volume"`
	High24Hrs             float64   `json:"high_24h"`
	Low24Hrs              float64   `json:"low_24h"`
	CirculatingSupply     float64   `json:"circulating_supply"`
	TotalSupply           float64   `json:"total_supply"`
	MaxSupply             float64   `json:"max_supply"`
	ATH                   float64   `json:"ath"`
	ATL                   float64   `json:"atl"`
	LastUpdated           time.Time `json:"last_updated"`
}

// Token contains the information of a single token
type Token struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

// Tokens represents a list of Token objects
type Tokens []Token

// TokenPrice represents the price of a token unit
type TokenPrice struct {
	CoingeckoID           string    `json:"id"`
	Symbol                string    `json:"symbol"`
	Name                  string    `json:"name"`
	Image                 string    `json:"image"`
	Price                 float64   `json:"price"`
	MarketCap             float64   `json:"market_cap"`
	MarketCapRank         int64     `json:"market_cap_rank"`
	FullyDilutedValuation float64   `json:"fully_diluted_valuation"`
	TotalVolume           float64   `json:"total_volume"`
	High24Hrs             float64   `json:"high_24h"`
	Low24Hrs              float64   `json:"low_24h"`
	CirculatingSupply     float64   `json:"circulating_supply"`
	TotalSupply           float64   `json:"total_supply"`
	MaxSupply             float64   `json:"max_supply"`
	ATH                   float64   `json:"ath"`
	ATL                   float64   `json:"atl"`
	Timestamp             time.Time `json:"timestamp"`
}

// NewTokenPrice creates new TokenPrice instance
func NewTokenPrice(coingeckoID, symbol, name, image string, price, marketCap float64, marketCapRank int64,
	fullyDilutedValuation, totalVolume, high24Hrs, low24Hrs, circulatingSupply, totalSupply,
	maxSupply, ath, atl float64, timestamp time.Time) TokenPrice {
	return TokenPrice{
		CoingeckoID:           coingeckoID,
		Symbol:                symbol,
		Name:                  name,
		Image:                 image,
		Price:                 price,
		MarketCap:             marketCap,
		MarketCapRank:         marketCapRank,
		FullyDilutedValuation: fullyDilutedValuation,
		TotalVolume:           totalVolume,
		High24Hrs:             high24Hrs,
		Low24Hrs:              low24Hrs,
		CirculatingSupply:     circulatingSupply,
		TotalSupply:           totalSupply,
		MaxSupply:             maxSupply,
		ATH:                   ath,
		ATL:                   atl,
		Timestamp:             timestamp,
	}
}
