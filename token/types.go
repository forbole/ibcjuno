package types

import "time"

type TokensConfig struct {
	Tokens []Token `yaml:"token"`
}

// Token represents valid token
type Token struct {
	Name  string      `yaml:"name"`
	Units []TokenUnit `yaml:"units"`
}

func NewToken(name string, units []TokenUnit) Token {
	return Token{
		Name:  name,
		Units: units,
	}
}

// TokenUnit represents a unit of a token
type TokenUnit struct {
	Denom    string `yaml:"denom"`
	IBCDenom string `yaml:"ibc_denom,omitempty"`
	Exponent int    `yaml:"exponent"`
	PriceID  string `yaml:"price_id,omitempty"`
}

func NewTokenUnit(denom string, ibcDenom string, exponent int, priceID string) TokenUnit {
	return TokenUnit{
		Denom:    denom,
		IBCDenom: ibcDenom,
		Exponent: exponent,
		PriceID:  priceID,
	}
}

// TokenPrice represents the price of a token unit
type TokenPrice struct {
	UnitName  string
	Price     float64
	MarketCap int64
	Timestamp time.Time
}

// NewTokenPrice returns TokenPrice instance
func NewTokenPrice(unitName string, price float64, marketCap int64, timestamp time.Time) TokenPrice {
	return TokenPrice{
		UnitName:  unitName,
		Price:     price,
		MarketCap: marketCap,
		Timestamp: timestamp,
	}
}

func NewTokensConfig(
	tokens []Token,
) TokensConfig {
	return TokensConfig{
		Tokens: tokens,
	}
}

// DefaultTokensConfig returns default TokensConfig instance
func DefaultTokensConfig() TokensConfig {
	var testToken []Token
	return NewTokensConfig(testToken)
}
