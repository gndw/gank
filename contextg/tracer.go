package contextg

import (
	"reflect"
	"runtime"
	"strings"
	"time"
)

type ContextGTracer struct {
	FunctionName string          `json:"function,omitempty"`
	T1           time.Time       `json:"t1,omitempty"`
	Time         time.Duration   `json:"time,omitempty"`
	Child        *ContextGTracer `json:"child,omitempty"`
}

func (t *ContextGTracer) Start(functionName string) {
	t.FunctionName = functionName
	t.T1 = time.Now()
}

func (t *ContextGTracer) Finish() {
	t.Time = time.Since(t.T1)
}

func FromFunction(i interface{}) string {
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	splits := strings.Split(fn, "/")
	if len(splits) > 0 {
		if splits[0] == "github.com" && len(splits) > 3 {
			return strings.Join(splits[3:], "/")
		}
	}
	return fn
}
