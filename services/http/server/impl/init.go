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
		return
	}

	port := 0

	// Get port from config file
	if params.Config != nil {
		port = params.Config.Server.Port
	}

	// Get port from Environment Machine
	if port <= 0 {
		portEnv := os.Getenv("PORT")
		if portEnvInt, err := strconv.Atoi(portEnv); err == nil && portEnvInt > 0 {
			port = portEnvInt
		}
	}

	// Set Default Port
	if port <= 0 {
		port = server.DEFAULT_PORT
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
	Config     *config.Service `optional:"true"`
	Lc         model.Lifecycle
	Shutdowner model.Shutdowner
	Router     router.Service
	Log        log.Service
}
