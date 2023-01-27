package global

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/config"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/video/videoservice"
	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig

	VideoClient videoservice.Client
)
