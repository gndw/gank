package impl

import (
	"context"
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

func New(params Parameters) (service server.Service, err error) {

	handler, err := params.Router.GetHandler()
	if err != nil {
		return nil, err
	}

	// Get port from config file
	port := params.Config.Server.Port

	// Replace port from Environment Machine if any
	portEnv := os.Getenv("PORT")
	if portEnvInt, err := strconv.Atoi(portEnv); err == nil && portEnvInt > 0 {
		port = portEnvInt
		params.Log.Debugf("server.service> using port %v found in machine environment variable", port)
	}

	serviceInstance := &Service{
		router: params.Router,
		server: &http.Server{Handler: handler},
	}

	params.Lc.Append(model.NewHook(model.Hook{
		OnStart: func(context.Context) (err error) {

			var listener net.Listener
			err = functions.LoggingProcessTime(params.Log, "starting http server", func() error {
				listener, err = net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
				return err
			})
			if err != nil {
				return err
			}

			params.Log.Infof("http server is listening at port %v", port)

			go func() {
				serveErr := serviceInstance.server.Serve(listener)
				if serveErr != nil && serveErr != http.ErrServerClosed {
					params.Log.Errorf("http server failed to start. err: ", err)
					params.Log.Warningf("shutting down application...")
					params.Shutdowner.Shutdown()
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return functions.LoggingProcessTime(params.Log, "stoping http server", func() error {
				return serviceInstance.server.Shutdown(ctx)
			})
		},
	}))

	return serviceInstance, nil
}

type Parameters struct {
	model.In
	Config     config.Service
	Lc         model.Lifecycle
	Shutdowner model.Shutdowner
	Router     router.Service
	Log        log.Service
}
