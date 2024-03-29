package api

type Config struct {
	ChainRegistryAPIURL string `yaml:"chain_registry_api_url"`
	ChainRegistryRawURL string `yaml:"chain_registry_raw_url"`
	CoingeckoURL        string `yaml:"coingecko_url"`
}

// NewAPIConfig creates new APIConfig instance
func NewAPIConfig(chainRegistryAPIURL, chainRegistryRawURL, coingeckoURL string) Config {
	return Config{
		ChainRegistryAPIURL: chainRegistryAPIURL,
		ChainRegistryRawURL: chainRegistryRawURL,
		CoingeckoURL:        coingeckoURL,
	}
}

// DefaultAPIConfig returns default APIConfig instance
func DefaultAPIConfig() Config {
	return NewAPIConfig(
		"https://api.github.com/repos/cosmos/chain-registry",
		"https://raw.githubusercontent.com/cosmos/chain-registry",
		"https://api.coingecko.com/api/v3")
}
