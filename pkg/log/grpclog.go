package log

import (
	"io/ioutil"
	"os"

	"google.golang.org/grpc/grpclog"
)

var Log grpclog.LoggerV2

func init() {
	Log = grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(Log)
}
