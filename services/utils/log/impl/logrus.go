package impl

import (
	"context"
	"strings"

	"github.com/gndw/gank/errorsg"
	"github.com/sirupsen/logrus"
)

func (s *Service) Print(args ...interface{}) {
	if len(args) > 0 {
		if arg, ok := args[0].(string); ok {
			if strings.Contains(arg, "ERROR") {
				s.Errorln(args...)
				return
			} else if strings.Contains(arg, "WARNING") {
				s.Warningln(args...)
				return
			} else if strings.Contains(arg, "INFO") {
				s.Infoln(args...)
				return
			} else if strings.Contains(arg, "DEBUG") {
				s.Debugln(args...)
				return
			} else if strings.Contains(arg, "TRACE") {
				s.Traceln(args...)
				return
			}
		}
	}
	s.Debugln(args...)
}

func (s *Service) Printf(str string, args ...interface{}) {
	if strings.Contains(str, "ERROR") {
		s.Errorf(str, args...)
		return
	}

	if len(args) > 0 {
		if arg, ok := args[0].(string); ok {
			if strings.Contains(arg, "ERROR") {
				s.Errorf(str, args...)
				return
			} else if strings.Contains(arg, "WARNING") {
				s.Warningf(str, args...)
				return
			} else if strings.Contains(arg, "INFO") {
				s.Infof(str, args...)
				return
			} else if strings.Contains(arg, "DEBUG") {
				s.Debugf(str, args...)
				return
			} else if strings.Contains(arg, "TRACE") {
				s.Tracef(str, args...)
				return
			}
		}
	}
	s.Debugf(str, args...)
}

func (s *Service) Tracef(str string, args ...interface{}) {
	logrus.Tracef(str, args...)
}

func (s *Service) Traceln(args ...interface{}) {
	logrus.Traceln(args...)
}

func (s *Service) Debugf(str string, args ...interface{}) {
	logrus.Debugf(str, args...)
}

func (s *Service) Debugln(args ...interface{}) {
	logrus.Debugln(args...)
}

func (s *Service) Infof(str string, args ...interface{}) {
	logrus.Infof(str, args...)
}

func (s *Service) Infoln(args ...interface{}) {
	logrus.Infoln(args...)
}

func (s *Service) Warningf(str string, args ...interface{}) {
	logrus.Warningf(str, args...)
}

func (s *Service) Warningln(args ...interface{}) {
	logrus.Warningln(args...)
}

func (s *Service) Errorf(str string, args ...interface{}) {
	logrus.Errorf(str, args...)
}

func (s *Service) Errorln(args ...interface{}) {
	logrus.Errorln(args...)
}

func (s *Service) InfoStd(ctx context.Context, msg string, metadata map[string]interface{}, err error) {

	stdMetadata := make(map[string]interface{})
	for key, value := range metadata {
		stdMetadata[key] = value
	}

	// get metadata from ctx
	// get metadata from error
	if err != nil {
		errMetadata := errorsg.GetMetadata(err)
		for key, value := range errMetadata {
			stdMetadata[key] = value
		}
	}

	logrus.WithFields(logrus.Fields(stdMetadata)).Infoln(msg)
}
