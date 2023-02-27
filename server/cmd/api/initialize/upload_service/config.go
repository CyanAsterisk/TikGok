package upload_service

import "github.com/CyanAsterisk/TikGok/server/cmd/api/config"

func initConfig() {
	conf = &config.GlobalServerConfig.UploadServiceInfo
}
