package global

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/user/config"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/sociality/socialityservice"
	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig

	SocialClient socialityservice.Client
)
