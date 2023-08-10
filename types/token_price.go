package types

import "time"

// MarketTicker contains the current market data for a single token
type MarketTicker struct {
	Symbol       string    `json:"symbol"`
	CurrentPrice float64   `json:"current_price"`
	MarketCap    float64   `json:"market_cap"`
	LastUpdated  time.Time `json:"last_updated"`
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
	UnitName  string
	Price     float64
	MarketCap int64
	Timestamp time.Time
}

// NewTokenPrice creates new TokenPrice instance
func NewTokenPrice(unitName string, price float64, marketCap int64, timestamp time.Time) TokenPrice {
	return TokenPrice{
		UnitName:  unitName,
		Price:     price,
		MarketCap: marketCap,
		Timestamp: timestamp,
	}
}
