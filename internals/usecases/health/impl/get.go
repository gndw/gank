package impl

import (
	"context"

	"github.com/gndw/gank/internals/constant"
	"github.com/gndw/gank/internals/model"
)

func (u *Usecase) Get(ctx context.Context) (result model.Health, err error) {

	result.Services = []model.ServiceHealth{}
	result.Services = append(result.Services, model.ServiceHealth{
		ServiceName: "app",
		IsHealthy:   true,
		Status:      constant.HealthOKResponse,
	})

	// check database health
	if u.db != nil {
		dberr := u.db.Ping(ctx)
		result.Services = append(result.Services, model.ServiceHealth{
			ServiceName: "db",
			IsHealthy:   dberr == nil,
			Status: func() string {
				if dberr != nil {
					return dberr.Error()
				}
				return constant.HealthOKResponse
			}(),
		})
	}

	return result, nil
}
