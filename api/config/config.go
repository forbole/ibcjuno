package api

type APIConfig struct {
	SupportedChainsURL string `yaml:"supported_chains_url"`
	ChainAssetsURL     string `yaml:"chain_assets_url"`
	CoingeckoURL       string `yaml:"coingecko_url"`
}

// NewAPIConfig creates new APIConfig instance
func NewAPIConfig(supportedChainsURL, chainAssetsURL, coingeckoURL string) APIConfig {
	return APIConfig{
		SupportedChainsURL: supportedChainsURL,
		ChainAssetsURL:     chainAssetsURL,
		CoingeckoURL:       coingeckoURL,
	}
}

// DefaultAPIConfig returns default APIConfig instance
func DefaultAPIConfig() APIConfig {
	return NewAPIConfig("", "", "")
}
