package config

type DatabaseConfig struct {
	Name               string `yaml:"name"`
	Host               string `yaml:"host"`
	Port               int64  `yaml:"port"`
	User               string `yaml:"user"`
	Password           string `yaml:"password"`
	SSLMode            string `yaml:"ssl_mode,omitempty"`
	Schema             string `yaml:"schema,omitempty"`
	MaxOpenConnections int    `yaml:"max_open_connections"`
	MaxIdleConnections int    `yaml:"max_idle_connections"`
}

// NewDatabaseConfig creates new DatabaseConfig instance
func NewDatabaseConfig(
	name, host string, port int64, user string, password string,
	sslMode string, schema string,
	maxOpenConnections int, maxIdleConnections int,
) DatabaseConfig {
	return DatabaseConfig{
		Name:               name,
		Host:               host,
		Port:               port,
		User:               user,
		Password:           password,
		SSLMode:            sslMode,
		Schema:             schema,
		MaxOpenConnections: maxOpenConnections,
		MaxIdleConnections: maxIdleConnections,
	}
}

// DefaultDatabaseConfig returns the default instance of DatabaseConfig
func DefaultDatabaseConfig() DatabaseConfig {
	return NewDatabaseConfig(
		"database-name",
		"localhost",
		5432,
		"user",
		"password",
		"",
		"public",
		1,
		1,
	)
}
