package utils

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ConfigParser represents a function that allows to parse config file
// contents as a Config object
type ConfigParser = func(fileContents []byte) (Config, error)

// DefaultConfigParser reads and parse IBCJuno config from the given string bytes.
// An error reading or parsing the config results in a panic.
func DefaultConfigParser(configData []byte) (Config, error) {
	var cfg = Config{
		bytes: configData,
	}
	err := yaml.Unmarshal(configData, &cfg)
	return cfg, err
}

// Read takes the path to a configuration file and returns the properly parsed configuration
func Read(configPath string, parser ConfigParser) (Config, error) {
	if configPath == "" {
		return Config{}, fmt.Errorf("empty configuration path")
	}

	configData, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config: %s", err)
	}

	return parser(configData)
}
