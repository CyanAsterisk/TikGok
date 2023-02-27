package initialize

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/api/config"
	config2 "github.com/CyanAsterisk/TikGok/server/cmd/api/pkg/uploadService/config"
)

func initConfig() {
	config2.GlobalServiceConfig = &config.GlobalServerConfig.UploadServiceInfo
}
