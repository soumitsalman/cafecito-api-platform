package shared

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

// Logging and error handling utilities
func NoError(err error, args ...any) {
	if err != nil {
		log.Fatal().Str("module", "DB").Err(err).Msg(fmt.Sprint(args...))
	}
}

func LogError(err error, msg string, args ...any) {
	if err != nil {
		log.Error().Str("module", "DB").Err(err).Msgf(msg, args...)
	}
}

func LogWarning(err error, msg string, args ...any) {
	if err != nil {
		log.Warn().Str("module", "DB").Err(err).Msgf(msg, args...)
	}
}

func LogResult(items any, err error) {
	if err != nil {
		log.Error().Str("module", "DB").Err(err).Msg("query failed")
	} else {
		log.Debug().Str("module", "DB").Interface("result", items).Msg("query succeeded")
	}
}

func LogQuery(query string, args map[string]any) {
	evt := log.Debug().Str("module", "DB").Str("sql", query)
	for key, value := range args {
		if key != "embedding" {
			evt = evt.Interface(key, value)
		} else {
			evt = evt.Str("embedding", "[REDACTED]") // Avoid logging large embeddings
		}
	}
	evt.Msg("query")
}
