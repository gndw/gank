package impl

import (
	"github.com/gndw/gank/services/config"
	"github.com/gndw/gank/services/flag"
	"github.com/gndw/gank/services/utils/log"
	"github.com/sirupsen/logrus"
)

type Service struct{}

func NewLogrus(flag flag.Service, config config.Service) (log.Service, error) {
	ins := &Service{}

	if config.Server.IsLoggingInJSON {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		})
	}

	if flag.Verbose != nil && *flag.Verbose {
		logrus.SetLevel(logrus.TraceLevel)
		ins.Debugln("log.service> verbose flag is active")
	}

	return ins, nil
}
