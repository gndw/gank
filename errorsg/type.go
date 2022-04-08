package errorsg

import (
	"encoding/json"
	"fmt"
	"strings"
)

type BuildOptions func(err CustomError) CustomError
type ErrorType string

var (
	ErrorTypeBadRequest          ErrorType = "bad_request"
	ErrorTypeInternalServerError ErrorType = "internal_server_error"
	ErrorTypePanic               ErrorType = "panic"
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

	// make copy of custom error to avoid altering original value
	ce := CustomError{

		Data:           e.Data,
		Type:           e.Type,
		HttpStatusCode: e.HttpStatusCode,
		Request:        e.Request,

		// Exclude private value that doesn't need to be printed
		// PrivateIdentifier: e.PrivateIdentifier,
		// PrettyMessage: e.PrettyMessage,
	}

	v, err := json.Marshal(ce)
	if err != nil {
		return "error object failed to marshal"
	}

	return string(v)
}

func GetMetadata(err error) (metadata map[string]interface{}) {

	metadata = make(map[string]interface{})
	metadata["error.data"] = GetData(err)

	exist, errorType := GetType(err)
	if exist {
		metadata["error.type"] = errorType
	}
	exist, httpStatusCode := GetHttpStatusCode(err)
	if exist {
		metadata["error.http_status_code"] = httpStatusCode
	}
	exist, request := GetRequest(err)
	if exist {
		for key, value := range request {
			bytevalue, _ := json.Marshal(value)
			metadata[fmt.Sprintf("error.request.%v", key)] = string(bytevalue)
		}
	}
	exist, pi := GetPrivateIdentifier(err)
	if exist {
		metadata["error.private_identifier"] = strings.Join(pi, "|")
	}
	exist, pmsg := GetPrettyMessage(err)
	if exist {
		metadata["error.pretty_message"] = pmsg
	}

	return metadata
}
