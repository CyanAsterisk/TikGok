package global

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/api/config"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/user/userservice"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/video/videoservice"
)

var (
	ServerConfig = &config.ServerConfig{}
	NacosConfig  = &config.NacosConfig{}

	UserClient  userservice.Client
	VideoClient videoservice.Client
)
