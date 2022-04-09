package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gndw/gank/errorsg"
	"github.com/gndw/gank/model"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func (s *Service) GetHttpMiddleware(f model.Middleware) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		var (
			data    interface{}
			err     error
			isPanic bool
		)

		ctx := r.Context()

		// if activate log
		t := time.Now()
		ww := middleware.NewWrapResponseWriter(rw, r.ProtoMajor)
		defer func() {
			s.LogHttpRequest(ctx, ww, r, time.Since(t), data, err)
		}()

		data, err = f(ctx, ww, r)

		response := model.HTTPResponse{}
		if data != nil {
			response.Data = data
		}
		if err != nil {

			isStatusCodeExist, statusCode := errorsg.GetHttpStatusCode(err)
			if isStatusCodeExist {
				render.Status(r, statusCode)
			} else {

				isErrorTypeExist, errorType := errorsg.GetType(err)
				if isErrorTypeExist {
					switch errorType {
					case errorsg.ErrorTypeBadRequest:
						render.Status(r, http.StatusBadRequest)
					case errorsg.ErrorTypeInternalServerError:
						render.Status(r, http.StatusInternalServerError)
					case errorsg.ErrorTypePanic:
						render.Status(r, http.StatusInternalServerError)
						isPanic = true
					}
				} else {
					// default unhandled error status
					render.Status(r, http.StatusInternalServerError)
				}
			}

			isPrettyMsgExist, prettyMsg := errorsg.GetPrettyMessage(err)
			if isPrettyMsgExist {
				response.Error = append(response.Error, prettyMsg)
			}

			response.Error = append(response.Error, err.Error())

		} else {
			render.Status(r, http.StatusOK)
		}

		if isPanic {
			// add custom content-length header if panic
			responseBytes, _ := json.Marshal(response)
			rw.Header().Add("Content-Length", strconv.Itoa(len(responseBytes)+1))
		}

		// write response
		render.JSON(ww, r, response)

		if isPanic {
			// flush if panic
			f, ok := rw.(http.Flusher)
			if ok {
				f.Flush()
			}
		}
	}
}

func (s *Service) LogHttpRequest(ctx context.Context, wrw middleware.WrapResponseWriter, r *http.Request, processTime time.Duration, data interface{}, err error) {

	stdMetadata := make(map[string]interface{})

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

	stdMetadata["process-time"] = processTime.Milliseconds()

	responseBytes, _ := json.Marshal(data)
	stdMetadata["response"] = string(responseBytes)

	// get metadata from ctx
	// get metadata from error
	if err != nil {
		errMetadata := errorsg.GetMetadata(err)
		for key, value := range errMetadata {
			stdMetadata[key] = value
		}
	}

	if err == nil {
		s.logService.LogInfoWithMetadata(stdMetadata, "http ok")
	} else {
		exist, errorType := errorsg.GetType(err)
		if exist {
			switch errorType {
			case errorsg.ErrorTypeBadRequest:
				s.logService.LogInfoWithMetadata(stdMetadata, "http bad request")
			case errorsg.ErrorTypeInternalServerError:
				s.logService.LogErrorWithMetadata(stdMetadata, "http internal server error")
			case errorsg.ErrorTypePanic:
				s.logService.LogPanicWithMetadata(stdMetadata, "http panic")
			default:
				s.logService.LogWarningWithMetadata(stdMetadata, "http unknown type")
			}
		} else {
			s.logService.LogWarningWithMetadata(stdMetadata, "http without type")
		}
	}
}
