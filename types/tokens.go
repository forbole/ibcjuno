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
	Denom    string         `yaml:"denom"`
	Exponent int            `yaml:"exponent"`
	PriceID  string         `yaml:"price_id,omitempty"`
	IBCDenom []IBCTokenUnit `yaml:"ibc_denom,omitempty"`
}

// IBCTokenUnit represents a unit of a IBC token
type IBCTokenUnit struct {
	Denom    string `yaml:"denom"`
	SrcChain string `yaml:"src_chain"`
	DstChain string `yaml:"dst_chain"`
	Channel  string `yaml:"channel"`
	IBCDenom string `yaml:"ibc_denom"`
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
	return NewTokensConfig([]Token{
		NewToken("Desmos", []TokenUnit{
			NewTokenUnit(
				"dsm",
				6,
				"desmos",
				[]IBCTokenUnit{
					NewIBCTokenUnit(
						"udsm",
						"desmos",
						"osmosis",
						"channel-1",
						"ibc/EA4C0A9F72E2CEDF10D0E7A9A6A22954DB3444910DB5BE980DF59B05A46DAD1C",
					),
				},
			),
		}),
	})
}

// NewToken creates new Token instance
func NewToken(name string, units []TokenUnit) Token {
	return Token{
		Name:  name,
		Units: units,
	}
}

// NewTokenUnit creates new TokenUnit instance
func NewTokenUnit(denom string, exponent int, priceID string, ibcDenom []IBCTokenUnit) TokenUnit {
	return TokenUnit{
		Denom:    denom,
		IBCDenom: ibcDenom,
		Exponent: exponent,
		PriceID:  priceID,
	}
}

// NewIBCTokenUnit creates new IBCTokenUnit instance
func NewIBCTokenUnit(denom string, srcChain string, dstChain string, channel string, ibcDenom string) IBCTokenUnit {
	return IBCTokenUnit{
		Denom:    denom,
		SrcChain: srcChain,
		DstChain: dstChain,
		Channel:  channel,
		IBCDenom: ibcDenom,
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
