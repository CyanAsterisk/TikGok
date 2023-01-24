package initialize

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzlogrus "github.com/hertz-contrib/obs-opentelemetry/logging/logrus"
)

// InitLogger to init logrus
func InitLogger() {
	logger := hertzlogrus.NewLogger()

	logger.SetLevel(hlog.LevelDebug)

	hlog.SetLogger(logger)
}
