package log

import "context"

type Service interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Tracef(string, ...interface{})
	Traceln(...interface{})
	Debugf(string, ...interface{})
	Debugln(...interface{})
	Infof(string, ...interface{})
	Infoln(...interface{})
	InfoStd(ctx context.Context, msg string, metadata map[string]interface{}, err error)
	Warningf(string, ...interface{})
	Warningln(...interface{})
	Errorf(string, ...interface{})
	Errorln(...interface{})
}
