package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gndw/gank/contextg"
	"github.com/gndw/gank/errorsg"
	"github.com/gndw/gank/model"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func (s *Service) GetHttpMiddleware(f model.Middleware) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		var (
			data      interface{}
			err       error
			errorType errorsg.ErrorType
		)

		// create custom context
		ctx := contextg.CreateCustomContext(r.Context())

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

			// calculate error type
			_, errorType = errorsg.GetType(err)

			// calculate http error code
			isStatusCodeExist, statusCode := errorsg.GetHttpStatusCode(err)
			if isStatusCodeExist {
				render.Status(r, statusCode)
			} else {
				// if no explicit http code provided, set error code based on error type
				switch errorType {
				case errorsg.ErrorTypeBadRequest:
					render.Status(r, http.StatusBadRequest)
				case errorsg.ErrorTypeInternalServerError:
					render.Status(r, http.StatusInternalServerError)
				case errorsg.ErrorTypePanic:
					render.Status(r, http.StatusInternalServerError)
				default:
					render.Status(r, http.StatusInternalServerError)
				}
			}

			// calculate error message
			isPrettyMsgExist, prettyMsg := errorsg.GetPrettyMessage(err)
			if isPrettyMsgExist {
				response.Error = append(response.Error, prettyMsg)
			} else {
				// if no explicit error message provided, set error message based on error type
				// also add request_id for better tracking in log engine
				_, reqID := contextg.GetRequestID(ctx)
				switch errorType {
				case errorsg.ErrorTypeBadRequest:
					response.Error = append(response.Error, fmt.Sprintf(s.configService.Server.DefaultMsgBadRequest, reqID))
				case errorsg.ErrorTypeInternalServerError:
					response.Error = append(response.Error, fmt.Sprintf(s.configService.Server.DefaultMsgInternalServerError, reqID))
				case errorsg.ErrorTypePanic:
					response.Error = append(response.Error, fmt.Sprintf(s.configService.Server.DefaultMsgInternalServerError, reqID))
				default:
					response.Error = append(response.Error, fmt.Sprintf(s.configService.Server.DefaultMsgInternalServerError, reqID))
				}
			}

			// adding secondary error message as actual error for developer
			if s.configService.Server.IsReturnDeveloperError {
				response.Error = append(response.Error, err.Error())
			}

		} else {

			// set status ok if no error provided
			render.Status(r, http.StatusOK)

		}

		// if error type is panic, set content-length to actual response length
		// to force golang-http sending error message back to the client
		if errorType == errorsg.ErrorTypePanic {
			responseBytes, _ := json.Marshal(response)
			rw.Header().Add("Content-Length", strconv.Itoa(len(responseBytes)+1))
		}

		// write response
		render.JSON(ww, r, response)

		// if error type is panic, flush after response is written
		// to force golang-http sending error message back to the client
		if errorType == errorsg.ErrorTypePanic {
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
