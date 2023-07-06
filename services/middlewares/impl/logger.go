package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gndw/gank/contextg"
	"github.com/gndw/gank/errorsg"
	"github.com/gndw/gank/functions"
	"github.com/gndw/gank/model"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Service) GetLoggerMiddleware(f model.Middleware, options ...model.MiddlewareOption) model.Middleware {
	return func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (data interface{}, err error) {

		// wrapper to catch response
		ww := middleware.NewWrapResponseWriter(rw, r.ProtoMajor)

		// next
		data, err = f(ctx, ww, r)

		// logging
		s.LogHttpRequest(ctx, ww, r, data, err, options...)

		return data, err
	}
}

var LoggerOptionAdditionalSensitiveFields model.MiddlewareOptionKey = "LoggerOptionAdditionalSensitiveFields"

func (s *Service) WithLoggerMiddlewareOption_AdditionalSensitiveFields(sensitiveFields []string) model.MiddlewareOption {
	return model.MiddlewareOption{
		Key:   LoggerOptionAdditionalSensitiveFields,
		Value: sensitiveFields,
	}
}

func (s *Service) LogHttpRequest(ctx context.Context, wrw middleware.WrapResponseWriter, r *http.Request, data interface{}, err error, options ...model.MiddlewareOption) {

	stdMetadata := make(map[string]interface{})

	// basic
	stdMetadata["env"] = s.envService.Get()
	stdMetadata["process-time"] = s.getProcessingTime(ctx)

	// http related logging
	stdMetadata["method"] = r.Method
	stdMetadata["endpoint"] = fmt.Sprintf("%s://%s%s %s\" ", s.getScheme(r), r.Host, r.RequestURI, r.Proto)
	stdMetadata["remote-address"] = r.RemoteAddr
	stdMetadata["status-code"] = wrw.Status()
	stdMetadata["bytes-written"] = wrw.BytesWritten()
	stdMetadata["returned-headers"] = s.getReturnedHeaders(wrw)

	// request body and response
	stdMetadata["request-body"] = s.getRequestBody(ctx, options...)
	stdMetadata["response"] = s.getResponse(data, options...)

	// get metadata from ctx
	ctxMetadata := contextg.GetMetadata(ctx)
	for key, value := range ctxMetadata {
		stdMetadata[key] = value
	}

	// get metadata from error
	if err != nil {
		errMetadata := errorsg.GetMetadata(err)
		for key, value := range errMetadata {
			stdMetadata[key] = value
		}
	}

	msg := fmt.Sprintf("HTTP Request | %v %v | code %v | in %v ms", stdMetadata["method"], stdMetadata["endpoint"], stdMetadata["status-code"], stdMetadata["process-time"])

	if err == nil {
		s.logService.LogInfoWithMetadata(stdMetadata, msg)
	} else {
		exist, errorType := errorsg.GetType(err)
		if exist {
			switch errorType {
			case errorsg.ErrorTypeBadRequest:
				s.logService.LogInfoWithMetadata(stdMetadata, msg)
			case errorsg.ErrorTypeInternalServerError:
				s.logService.LogErrorWithMetadata(stdMetadata, msg)
			case errorsg.ErrorTypePanic:
				s.logService.LogPanicWithMetadata(stdMetadata, msg)
			default:
				s.logService.LogWarningWithMetadata(stdMetadata, msg)
			}
		} else {
			s.logService.LogWarningWithMetadata(stdMetadata, msg)
		}
	}
}

func (s *Service) getScheme(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme
}

func (s *Service) getReturnedHeaders(wrw middleware.WrapResponseWriter) string {
	returnedHeadersBytes, _ := json.Marshal(wrw.Header())
	return string(returnedHeadersBytes)
}

func (s *Service) getProcessingTime(ctx context.Context) string {
	exist, t1 := contextg.GetIncomingTime(ctx)
	if exist {
		return fmt.Sprint(time.Since(t1).Milliseconds())
	}
	return "-"
}

func (s *Service) getRequestBody(ctx context.Context, options ...model.MiddlewareOption) string {
	exist, rb := contextg.GetRequestBody(ctx)
	if exist {
		return string(functions.MaskingDataFromBytes(rb, s.getSensitiveFields(options...)))
	}
	return ""
}

func (s *Service) getResponse(data interface{}, options ...model.MiddlewareOption) string {
	responseBytes, _ := json.Marshal(data)
	return string(functions.MaskingDataFromBytes(responseBytes, s.getSensitiveFields(options...)))
}

func (s *Service) getSensitiveFields(options ...model.MiddlewareOption) []string {
	sensitiveFields := strings.Split(s.configService.Server.SensitiveFields, ",")
	for _, option := range options {
		if option.Key == LoggerOptionAdditionalSensitiveFields {
			if additionalSensitiveFields, ok := option.Value.([]string); ok {
				sensitiveFields = append(sensitiveFields, additionalSensitiveFields...)
			}
		}
	}
	fmt.Println("sf", sensitiveFields)
	return sensitiveFields
}
