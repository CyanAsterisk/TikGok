package initialize

import (
	"github.com/cloudwego/kitex/pkg/klog"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
)

// InitLogger to init logrus
func InitLogger() {
	logger := kitexlogrus.NewLogger()

	logger.SetLevel(klog.LevelDebug)

	klog.SetLogger(logger)
}
