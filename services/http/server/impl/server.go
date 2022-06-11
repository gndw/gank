package impl

import (
	"github.com/gndw/gank/model"
)

func (s *Service) AddHttpHandler(req model.AddHTTPRequest) (err error) {
	return s.router.AddHttpHandler(req)
}

func (s *Service) AddHttpHandlers(requests ...model.AddHTTPRequest) (err error) {
	for _, request := range requests {
		err = s.router.AddHttpHandler(request)
		if err != nil {
			return err
		}
	}
	return nil
}
