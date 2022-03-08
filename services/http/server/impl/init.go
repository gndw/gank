package impl

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/gndw/gank/functions"
	"github.com/gndw/gank/model"
	"github.com/gndw/gank/services/config"
	"github.com/gndw/gank/services/http/router"
	"github.com/gndw/gank/services/http/server"
	"github.com/gndw/gank/services/utils/log"
)

type Service struct {
	server *http.Server
	router router.Service
}

func New(config config.Service, lc model.Lifecycle, shutdowner model.Shutdowner, router router.Service, log log.Service) (service server.Service, err error) {

	handler, err := router.GetHandler()
	if err != nil {
		return
	}

	port := config.GetServer().GetPort()

	portEnv := os.Getenv("PORT")
	if portEnvInt, err := strconv.Atoi(portEnv); err == nil && portEnvInt > 0 {
		port = portEnvInt
	}

	if port <= 0 {
		return nil, errors.New("server port not found in config")
	}

	serviceInstance := &Service{
		router: router,
		server: &http.Server{Handler: handler},
	}

	lc.Append(model.NewHook(model.Hook{
		OnStart: func(context.Context) (err error) {

			var listener net.Listener
			err = functions.LoggingProcessTime(log, "starting http server", func() error {
				listener, err = net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
				return err
			})
			if err != nil {
				return err
			}

			log.Infof("http server is listening at port %v", port)

			go func() {
				serveErr := serviceInstance.server.Serve(listener)
				if serveErr != nil && serveErr != http.ErrServerClosed {
					log.Errorf("http server failed to start. err: ", err)
					log.Warningf("shutting down application...")
					shutdowner.Shutdown()
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return functions.LoggingProcessTime(log, "stoping http server", func() error {
				return serviceInstance.server.Shutdown(ctx)
			})
		},
	}))

	return serviceInstance, nil
}
