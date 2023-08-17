package api

type APIConfig struct {
	ChainRegistryURL       string `yaml:"chain_registry_url"`
	ChainRegistryAssetsURL string `yaml:"chain_registry_assets_url"`
	CoingeckoURL           string `yaml:"coingecko_url"`
}

// NewAPIConfig creates new APIConfig instance
func NewAPIConfig(chainRegistryURL, chainRegistryAssetsURL, coingeckoURL string) APIConfig {
	return APIConfig{
		ChainRegistryURL:       chainRegistryURL,
		ChainRegistryAssetsURL: chainRegistryAssetsURL,
		CoingeckoURL:           coingeckoURL,
	}
}

// DefaultAPIConfig returns default APIConfig instance
func DefaultAPIConfig() APIConfig {
	return NewAPIConfig("", "", "")
}
