package config

import (
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

// GetEnv returns the trimmed value of name, then fallback if provided.
// It terminates the process when a required value is missing.
func GetEnv(name, fallback string, mustExist bool) string {
	if v := strings.TrimSpace(os.Getenv(name)); v != "" {
		return v
	} else if fallback != "" {
		return fallback
	}
	if mustExist {
		log.Fatal().Str("module", "MAIN").Msgf("%s is required", name)
	}
	return fallback
}

// ParseAPIKeys converts a semicolon-separated Header=Value string to API keys.
func ParseAPIKeys(raw string) map[string]string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	result := map[string]string{}
	for _, pair := range strings.Split(raw, ";") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			continue
		}
		header := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if header != "" && value != "" {
			result[header] = value
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}
