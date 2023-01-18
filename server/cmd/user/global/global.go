package global

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/user/config"
	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig
)
