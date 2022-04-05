package errorsg

import (
	"encoding/json"
)

type BuildOptions func(err CustomError) CustomError
type ErrorType string

var (
	ErrorTypeBadRequest          ErrorType = "bad_request"
	ErrorTypeInternalServerError ErrorType = "internal_server_error"
)

type CustomError struct {
	Data              string                  `json:"data,omitempty"`
	Type              *ErrorType              `json:"type,omitempty"`
	HttpStatusCode    *int                    `json:"http_status_code,omitempty"`
	Request           *map[string]interface{} `json:"request,omitempty"`
	PrivateIdentifier *[]string               `json:"private_identifier,omitempty"`
	PrettyMessage     *string                 `json:"pretty_massage,omitempty"`
}

func (e *CustomError) Error() string {

	// Exclude any private field before printing error
	e.PrivateIdentifier = nil

	// Exclude pretty message before printing error
	e.PrettyMessage = nil

	v, err := json.Marshal(e)
	if err != nil {
		return "error object failed to marshal"
	}

	return string(v)
}
