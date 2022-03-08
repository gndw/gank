package impl

import (
	"context"
	"net/http"

	"github.com/gndw/gank/errorsg"
	"github.com/gndw/gank/model"
	"github.com/go-chi/render"
)

func (s *Service) GetHttpMiddleware(f model.Middleware) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		ctx := context.Background()

		data, err := f(ctx, rw, r)

		response := model.HTTPResponse{}
		if data != nil {
			response.Data = data
		}
		if err != nil {

			isStatusCodeExist, statusCode := errorsg.GetStatusCode(err)
			if isStatusCodeExist {
				render.Status(r, statusCode)
			} else {
				render.Status(r, 400)
			}

			isPrettyMsgExist, prettyMsg := errorsg.GetPrettyMessage(err)
			if isPrettyMsgExist {
				response.Error = append(response.Error, prettyMsg)
			}

			response.Error = append(response.Error, err.Error())
		}

		render.JSON(rw, r, response)
	}
}
