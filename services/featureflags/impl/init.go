package impl

import (
	"errors"
	"fmt"
	"io/ioutil"

	"net/http"
	"time"

	"github.com/gndw/gank/functions"
	"github.com/gndw/gank/services/config"
	"github.com/gndw/gank/services/featureflags"
	"github.com/gndw/gank/services/utils/log"
	"github.com/gndw/gank/services/utils/marshal"
)

type FeatureFlagKeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Service struct {
	FeatureFlags map[string]string
}

func NewRemoteJsonURL(config config.Service, log log.Service, marshal marshal.Service) (featureflags.Service, error) {
	ins := &Service{
		FeatureFlags: map[string]string{},
	}

	if config.FeatureFlags.RemoteFeatureFlagsJsonURL == "" {
		return ins, errors.New("cannot read remoteFeatureFlagsJsonURL in config file")
	}

	if config.FeatureFlags.WatchTimeInSeconds <= 0 {
		return ins, fmt.Errorf("feature flags watch time is invalid: %v", config.FeatureFlags.WatchTimeInSeconds)
	}

	readFeatureFlag := func() error {

		var readErr error
		var content []byte

		if functions.IsUrl(config.FeatureFlags.RemoteFeatureFlagsJsonURL) {

			// assume its remote URL
			resp, readErr := http.Get(config.FeatureFlags.RemoteFeatureFlagsJsonURL)
			if readErr != nil {
				return fmt.Errorf("error when fetching feature flag URL. err:%v", readErr)
			}

			// TODO create utility function to do http request
			content, readErr = ioutil.ReadAll(resp.Body)
			if readErr != nil {
				return fmt.Errorf("error when reading response body. err:%v", readErr)
			}

		} else {

			// assume its file path
			filepath, readErr := functions.GetPathFromDirtyPath(config.FeatureFlags.RemoteFeatureFlagsJsonURL)
			if readErr != nil {
				return fmt.Errorf("error when checking & sanitizing file path. err:%v", readErr)
			}

			content, readErr = ioutil.ReadFile(filepath)
			if readErr != nil {
				return fmt.Errorf("error when reading file. err:%v", readErr)
			}

		}

		featureFlags := []FeatureFlagKeyValuePair{}
		readErr = marshal.JsonUnmarshal(content, &featureFlags)
		if readErr != nil {
			return fmt.Errorf("error when converting response body to json. err:%v", readErr)
		}

		ins.FeatureFlags = make(map[string]string, len(featureFlags))
		for _, kv := range featureFlags {
			ins.FeatureFlags[kv.Key] = kv.Value
		}

		return nil
	}

	// watching feature flag
	go func() {
		// TODO: shutdown gracefully?
		for {
			log.Debugln("reading feature flag...")
			readErr := readFeatureFlag()
			if readErr != nil {
				log.Errorf(readErr.Error())
			}
			time.Sleep(time.Duration(config.FeatureFlags.WatchTimeInSeconds) * time.Second)
		}
	}()

	return ins, nil
}
