package errorsg

import (
	"fmt"
)

func WithOptions(err error, options ...BuildOptions) error {

	obj, ok := err.(*CustomError)
	if !ok {
		obj = &CustomError{
			Data: err.Error(),
		}
	}

	for _, opt := range options {
		p := opt(*obj)
		obj = &p
	}

	return obj
}

func BadRequestWithOptions(err error, options ...BuildOptions) error {
	options = append([]BuildOptions{
		WithType(ErrorTypeBadRequest),
	}, options...)
	return WithOptions(err, options...)
}

func InternalErrorWithOptions(err error, options ...BuildOptions) error {
	options = append([]BuildOptions{
		WithType(ErrorTypeInternalServerError),
	}, options...)
	return WithOptions(err, options...)
}

func JSONParseErrorWithOptions(err error, options ...BuildOptions) error {
	return BadRequestWithOptions(fmt.Errorf("failed to parse json string. err: %v", err.Error()), options...)
}
