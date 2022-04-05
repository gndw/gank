package impl

import (
	"net/http"

	"github.com/gndw/gank/errorsg"
	"github.com/gndw/gank/model"
	"github.com/go-chi/render"
)

func (s *Service) GetHttpMiddleware(f model.Middleware) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		data, err := f(ctx, rw, r)

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
		}

		s.logService.InfoStd(ctx, "incoming-http", nil, err)
		render.JSON(rw, r, response)
	}
}
