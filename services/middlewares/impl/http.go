package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gndw/gank/contextg"
	"github.com/gndw/gank/errorsg"
	"github.com/gndw/gank/model"
	"github.com/go-chi/render"
)

func (s *Service) GetHttpMiddleware(f model.Middleware, options ...model.MiddlewareOption) model.Middleware {
	return func(ctx context.Context, rw http.ResponseWriter, r *http.Request) (data interface{}, err error) {

		var (
			errorType errorsg.ErrorType
		)

		// catch request body
		b, e := io.ReadAll(r.Body)
		if e == nil {
			ctx = contextg.WithRequestBody(ctx, b)
		}
		r.Body.Close()

		// execute next
		data, err = f(ctx, rw, r)

		// calculate response
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
		render.JSON(rw, r, response)

		// if error type is panic, flush after response is written
		// to force golang-http sending error message back to the client
		if errorType == errorsg.ErrorTypePanic {
			f, ok := rw.(http.Flusher)
			if ok {
				f.Flush()
			}
		}

		return data, err
	}
}
