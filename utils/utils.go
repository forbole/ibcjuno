package utils

import (
	"database/sql"
	"os"
	"path"
	"strings"

	"github.com/rs/zerolog/log"
	yaml "gopkg.in/yaml.v3"
)

var (
	HomePath = ""
)

// WatchMethod allows to watch for a method that returns an error.
// It executes the given method in a goroutine, logging any error that might raise.
func WatchMethod(method func() error) {
	go func() {
		err := method()
		if err != nil {
			log.Error().Err(err)
		}
	}()
}

// ToNullString converts to empty string
func ToNullString(value string) sql.NullString {
	value = strings.TrimSpace(value)
	return sql.NullString{
		Valid:  value != "",
		String: value,
	}
}

// GetConfigFilePath returns the path to the configuration file given the executable name
func GetConfigFilePath() string {
	return path.Join(HomePath, "config.yaml")
}

// Write allows to write the given configuration into the file present at the given path
func Write(cfg Config, path string) error {
	bz, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, bz, 0600)
}

// RemoveDuplicates function removes duplicate values
func RemoveDuplicates(s []string) []string {
	bucket := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := bucket[str]; !ok {
			bucket[str] = true
			result = append(result, str)
		}
	}
	return result
}
