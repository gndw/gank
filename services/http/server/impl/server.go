package impl

import (
	"github.com/gndw/gank/model"
)

func (s *Service) AddHttpHandler(req model.AddHTTPRequest) (err error) {
	return s.router.AddHttpHandler(req)
}
