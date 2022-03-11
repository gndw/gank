package impl

import (
	"github.com/gndw/gank/services/flag"
	"github.com/gndw/gank/services/utils/log"
	"github.com/sirupsen/logrus"
)

type Service struct{}

func NewLogrus(flag flag.Service) (log.Service, error) {
	ins := &Service{}

	// TODO: this is dev only ya
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	if flag.Verbose != nil && *flag.Verbose {
		logrus.SetLevel(logrus.TraceLevel)
	}

	return ins, nil
}
