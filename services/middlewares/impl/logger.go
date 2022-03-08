package impl

import (
	"net/http"
	"runtime"

	"github.com/gndw/gank/services/utils/log"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Service) GetLoggerMiddleware() func(next http.Handler) http.Handler {

	if s.logMiddleware == nil {
		color := true
		if runtime.GOOS == "windows" {
			color = false
		}
		customLogger := &CustomHTTPLogger{
			log: s.logService,
		}
		s.logMiddleware = middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: customLogger, NoColor: !color})
	}

	return s.logMiddleware
}

type CustomHTTPLogger struct {
	log log.Service
}

func (c *CustomHTTPLogger) Print(v ...interface{}) {
	v = append([]interface{}{"INFO"}, v...)
	c.log.Print(v...)
}
