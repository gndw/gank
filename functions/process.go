package functions

import (
	"time"

	"github.com/gndw/gank/services/utils/log"
)

func LoggingProcessTime(log log.Service, processName string, process func() error, opts ...LoggingProcessTimeOption) error {

	cfg := LoggingProcessTimeConfig{}
	for _, opt := range opts {
		cfg = opt(cfg)
	}

	log.Infof("s> %v %v", processName, "...")

	startTime := time.Now()
	err := process()
	if err != nil {
		log.Errorf("FAILED> %v in %v", processName, time.Since(startTime))
	} else {
		t := time.Since(startTime)
		if cfg.WarningTimeLimitInSeconds != nil && t.Seconds() > *cfg.WarningTimeLimitInSeconds {
			log.Warningf("d> %v in %v", processName, time.Since(startTime))
		} else {
			log.Infof("d> %v in %v", processName, time.Since(startTime))
		}
	}

	return err
}

type LoggingProcessTimeConfig struct {
	WarningTimeLimitInSeconds *float64
}

type LoggingProcessTimeOption func(cfg LoggingProcessTimeConfig) LoggingProcessTimeConfig

func WithLoggingProcessTimeLimit(seconds float64) LoggingProcessTimeOption {
	return func(cfg LoggingProcessTimeConfig) LoggingProcessTimeConfig {
		cfg.WarningTimeLimitInSeconds = &seconds
		return cfg
	}
}
