package config

type ParserConfig struct {
	RefreshIBCTokensOnStart bool `yaml:"refresh_ibc_tokens_on_start"`
}

// NewParserConfig creates new ParserConfig instance
func NewParserConfig(
	refreshIBCTokensOnStart bool,
) ParserConfig {
	return ParserConfig{
		RefreshIBCTokensOnStart: refreshIBCTokensOnStart,
	}
}

// DefaultParserConfig returns the default instance of ParserConfig
func DefaultParserConfig() ParserConfig {
	return NewParserConfig(true)

}
