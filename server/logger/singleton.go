package logger

import (
	"os"
	"sync"

	"github.com/go-kit/kit/log"
)

var instance *log.Logger
var once sync.Once

//GetLogger get singleton logger
func GetLogger() *log.Logger {
	once.Do(func() {
		logger := log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
		instance = &logger
	})
	return instance
}
