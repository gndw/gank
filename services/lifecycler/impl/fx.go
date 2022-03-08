package impl

import (
	"context"
	"fmt"
	golog "log"
	"runtime/debug"
	"time"

	"github.com/gndw/gank/model"
	"github.com/gndw/gank/services/utils/log"
	"go.uber.org/dig"
	"go.uber.org/fx"
)

func (s *Service) Run() (err error) {
	err = s.app.Start(context.Background())
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("application startup in %v", time.Since(s.startTime))
	if s.log != nil {
		s.log.Infoln(msg)
	} else {
		golog.Println(msg)
	}

	<-s.app.Done()

	fmt.Println(" >> exiting...")
	doneTime := time.Now()

	err = s.app.Stop(context.Background())
	if err != nil {
		return err
	}

	msg = fmt.Sprintf("application exit in %v", time.Since(doneTime))
	if s.log != nil {
		s.log.Infoln(msg)
	} else {
		golog.Println(msg)
	}
	return nil
}

func (s *Service) PreConfig() (err error) {
	s.startTime = time.Now()
	s.container = dig.New()
	s.fxOptions = append(s.fxOptions, fx.Invoke(func(lc fx.Lifecycle, sd fx.Shutdowner) {
		s.container.Provide(func() model.Lifecycle {
			my := model.MyLifecycle{
				Fxl:            lc,
				OnInvokeErrors: []func(context.Context) error{},
			}
			return &my
		})
		s.container.Provide(func() model.Shutdowner {
			return sd
		})
	}))

	return nil
}

func (s *Service) PostConfig() (err error) {
	s.app = fx.New(s.fxOptions...)
	return s.app.Err()
}

func (s *Service) AddProviders(providers ...interface{}) (err error) {
	for _, provider := range providers {
		err = s.container.Provide(provider)
		if err != nil {
			return
		}
	}
	return nil
}
func (s *Service) AddInvokers(invokers ...interface{}) (err error) {
	for i := range invokers {
		index := i
		s.fxOptions = append(s.fxOptions, fx.Invoke(func() (err error) {

			// recovery and shutdown when panic
			defer func() {
				if r := recover(); r != nil {
					Shutdown(s.container)
					err = fmt.Errorf("panic at startup! : %s \n %s", r, string(debug.Stack()))
				}
			}()

			// shutdown when invoke / startup is failed
			err = s.container.Invoke(invokers[index])
			if err != nil {
				Shutdown(s.container)
				return err
			}

			return nil
		}))
	}
	return nil
}

func Shutdown(dig *dig.Container) {
	dig.Invoke(func(applog log.Service, lc model.Lifecycle, shutdowner model.Shutdowner) {
		if applog != nil {
			applog.Warningf("shutting down application...")
		}
		if shutdowner != nil {
			shutdowner.Shutdown()
		}
		if lc != nil {
			lc.ExecuteOnErrors(context.Background())
		}
	})
}
