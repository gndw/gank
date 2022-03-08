package impl

import (
	"github.com/gndw/gank/services/utils/log"
	"github.com/sirupsen/logrus"
)

type Service struct{}

func NewLogrus() (log.Service, error) {
	ins := &Service{}

	// TODO: this is dev only ya
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	// // TODO: this is dev only ya
	// logrus.SetLevel(logrus.DebugLevel)

	return ins, nil
}
