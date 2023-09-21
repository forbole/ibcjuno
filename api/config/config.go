package api

type APIConfig struct {
	ChainRegistryURL string `yaml:"chain_registry_url"`
	CoingeckoURL     string `yaml:"coingecko_url"`
}

// NewAPIConfig creates new APIConfig instance
func NewAPIConfig(chainRegistryURL, coingeckoURL string) APIConfig {
	return APIConfig{
		ChainRegistryURL: chainRegistryURL,
		CoingeckoURL:     coingeckoURL,
	}
}

// DefaultAPIConfig returns default APIConfig instance
func DefaultAPIConfig() APIConfig {
	return NewAPIConfig("", "")
}
