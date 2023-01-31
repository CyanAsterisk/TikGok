package global

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/video/config"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction/interactionserver"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/user/userservice"
	"gorm.io/gorm"
)

var (
	DB             *gorm.DB
	ServerConfig   config.ServerConfig
	NacosConfig    config.NacosConfig
	InteractClient interactionserver.Client
	UserClient     userservice.Client
)
