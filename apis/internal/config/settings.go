package config

import "time"

const (
	VECTOR_QUERY_DEFAULT_CANDIDATE_LIMIT    = 512
	VECTOR_QUERY_CANDIDATE_LIMIT_MULTIPLIER = 4
	RETRY_ATTEMPTS                          = 5
	RETRY_DELAY                             = 30 * time.Second
)
