package global

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/video/config"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction/interactionserver"
	"gorm.io/gorm"
)

var (
	DB             *gorm.DB
	ServerConfig   config.ServerConfig
	NacosConfig    config.NacosConfig
	InteractClient interactionserver.Client
)
