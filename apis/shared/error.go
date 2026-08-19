package shared

import "fmt"

const (
	API_ERROR_INVALID_REQUEST = "invalid_request"
	API_ERROR_DB_ERROR        = "db_unavailable"
	API_ERROR_EMBEDDING_ERROR = "embedder_unavailable"
	API_ERROR_ENCODING_ERROR  = "encoding_error"
	API_ERROR_INVALID_DATA    = "invalid_data"
	API_ERROR_NOT_FOUND       = "not_found"
	API_ERROR_UNAUTHORIZED    = "unauthorized"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewAPIError(code string, message string) APIError {
	return APIError{Code: code, Message: message}
}

func (e APIError) Error() string {
	return fmt.Sprintf("[%s]: %s", e.Code, e.Message)
}
