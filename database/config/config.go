package config

type DatabaseConfig struct {
	URL                string `yaml:"url"`
	MaxOpenConnections int    `yaml:"max_open_connections"`
	MaxIdleConnections int    `yaml:"max_idle_connections"`
	SSLModeEnable      string `yaml:"ssl_mode_enable"`
	SSLRootCert        string `yaml:"ssl_root_cert"`
	SSLCert            string `yaml:"ssl_cert"`
	SSLKey             string `yaml:"ssl_key"`
}

// NewDatabaseConfig creates new DatabaseConfig instance
func NewDatabaseConfig(
	url, sslModeEnable, sslRootCert, sslCert, sslKey string,
	maxOpenConnections int, maxIdleConnections int,
) DatabaseConfig {
	return DatabaseConfig{
		URL:                url,
		MaxOpenConnections: maxOpenConnections,
		MaxIdleConnections: maxIdleConnections,
		SSLModeEnable:      sslModeEnable,
		SSLRootCert:        sslRootCert,
		SSLCert:            sslCert,
		SSLKey:             sslKey,
	}
}

// DefaultDatabaseConfig returns the default instance of DatabaseConfig
func DefaultDatabaseConfig() DatabaseConfig {
	return NewDatabaseConfig(
		"postgresql://user:password@localhost:5432/database-name?sslmode=disable&search_path=public",
		"false",
		"",
		"",
		"",
		1,
		1,
	)
}
