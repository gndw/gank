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
	"github.com/gndw/gank/model"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Service) GetLoggerMiddleware(f model.Middleware) model.Middleware {
	return func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (data interface{}, err error) {

		// wrapper to catch response
		ww := middleware.NewWrapResponseWriter(rw, r.ProtoMajor)

		// next
		data, err = f(ctx, ww, r)

		// logging
		s.LogHttpRequest(ctx, ww, r, data, err)

		return data, err
	}
}

func (s *Service) LogHttpRequest(ctx context.Context, wrw middleware.WrapResponseWriter, r *http.Request, data interface{}, err error) {

	stdMetadata := make(map[string]interface{})

	stdMetadata["env"] = s.envService.Get()

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	stdMetadata["method"] = r.Method
	stdMetadata["endpoint"] = fmt.Sprintf("%s://%s%s %s\" ", scheme, r.Host, r.RequestURI, r.Proto)
	stdMetadata["remote-address"] = r.RemoteAddr

	stdMetadata["status-code"] = wrw.Status()
	stdMetadata["bytes-written"] = wrw.BytesWritten()

	returnedHeadersBytes, _ := json.Marshal(wrw.Header())
	stdMetadata["returned-headers"] = string(returnedHeadersBytes)

	exist, t1 := contextg.GetIncomingTime(ctx)
	if exist {
		stdMetadata["process-time"] = time.Since(t1).Milliseconds()
	}

	responseBytes, _ := json.Marshal(data)
	stdMetadata["response"] = string(s.SanitizeSensitiveDataFromBytes(responseBytes))

	exist, rb := contextg.GetRequestBody(ctx)
	if exist {
		stdMetadata["request-body"] = string(s.SanitizeSensitiveDataFromBytes(rb))
	}

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

func (s *Service) SanitizeSensitiveDataFromBytes(from []byte) (to []byte) {

	var f interface{}
	json.Unmarshal(from, &f)

	v, ok := f.(map[string]interface{})
	if ok {

		sfields := strings.Split(s.configService.Server.SensitiveFields, ",")
		for _, sfield := range sfields {

			layers := strings.Split(sfield, ".")
			temp := &v

			for i, layer := range layers {

				if i == len(layers)-1 {
					_, exist := (*temp)[layer]
					if exist {
						(*temp)[layer] = "-MASKED-"
					}
				} else {
					_, exist := (*temp)[layer]
					if exist {
						temp2, okk := (*temp)[layer].(map[string]interface{})
						if okk {
							temp = &temp2
						} else {
							break
						}
					} else {
						break
					}
				}
			}
		}
	}

	rb, _ := json.Marshal(v)
	return rb
}
