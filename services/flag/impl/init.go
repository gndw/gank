package impl

import (
	goflag "flag"

	"github.com/gndw/gank/services/flag"
)

func New() (flag.Service, error) {

	ins := flag.Service{}

	ins.Env = goflag.String("env", "", "process environment")
	ins.Verbose = goflag.Bool("v", false, "verbose")
	goflag.Parse()

	return ins, nil
}
