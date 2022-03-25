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

// TokenUnit represents a unit of a token
type TokenUnit struct {
	Denom    string `yaml:"denom"`
	IBCDenom string `yaml:"ibc_denom,omitempty"`
	Exponent int    `yaml:"exponent"`
	PriceID  string `yaml:"price_id,omitempty"`
}

// TokenPrice represents the price of a token unit
type TokenPrice struct {
	UnitName  string
	Price     float64
	MarketCap int64
	Timestamp time.Time
}

// NewTokensConfig creates new TokensConfig instance
func NewTokensConfig(tokens []Token) TokensConfig {
	return TokensConfig{
		Tokens: tokens,
	}
}

// DefaultTokensConfig returns default TokensConfig instance
func DefaultTokensConfig() TokensConfig {
	var testToken []Token
	return NewTokensConfig(testToken)
}

// NewToken creates new Token instance
func NewToken(name string, units []TokenUnit) Token {
	return Token{
		Name:  name,
		Units: units,
	}
}

// NewTokenUnit creates new TokenUnit instance
func NewTokenUnit(denom string, ibcDenom string, exponent int, priceID string) TokenUnit {
	return TokenUnit{
		Denom:    denom,
		IBCDenom: ibcDenom,
		Exponent: exponent,
		PriceID:  priceID,
	}
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
