package types

type CounterPartyDetails struct {
	Channel string `yaml:"channel"`
	Port    string `yaml:"port"`
	Denom   string `yaml:"denom"`
}

// IBCTokenUnit represents a unit of a IBC token
type IBCTokenUnit struct {
	Denom        string              `json:"denom" yaml:"denom"`
	OriginChain  string              `json:"origin_chain" yaml:"origin_chain"`
	OriginDenom  string              `json:"origin_denom" yaml:"origin_denom"`
	OriginType   string              `json:"origin_type" yaml:"origin_type"`
	Symbol       string              `json:"symbol" yaml:"symbol"`
	Decimals     int64               `json:"decimals" yaml:"decimals"`
	Enable       bool                `json:"enable" yaml:"enable"`
	Path         string              `json:"path" yaml:"path"`
	Channel      string              `json:"channel" yaml:"channel"`
	CounterParty CounterPartyDetails `json:"counter_party" yaml:"counter_party"`
	CoingeckoID  string              `json:"coinGeckoId" yaml:"coinGeckoId"`
}

// NewIBCTokenUnit creates new IBCTokenUnit instance
func NewIBCTokenUnit(denom string, originChain, originDenom, originType, symbol string,
	decimals int64, enable bool, path, channel string,
	counterParty CounterPartyDetails, coingeckoID string) IBCTokenUnit {
	return IBCTokenUnit{
		Denom:        denom,
		OriginChain:  originChain,
		OriginDenom:  originDenom,
		OriginType:   originType,
		Symbol:       symbol,
		Decimals:     decimals,
		Enable:       enable,
		Path:         path,
		Channel:      channel,
		CounterParty: counterParty,
		CoingeckoID:  coingeckoID,
	}
}
