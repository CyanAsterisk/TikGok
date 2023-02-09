package global

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/api/config"
	"github.com/CyanAsterisk/TikGok/server/cmd/api/pkg/uploadService"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/chat/chatservice"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction/interactionserver"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/sociality/socialityservice"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/user/userservice"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/video/videoservice"
)

var (
	ServerConfig  = &config.ServerConfig{}
	NacosConfig   = &config.NacosConfig{}
	UploadService *uploadService.Service

	ChatClient        chatservice.Client
	UserClient        userservice.Client
	VideoClient       videoservice.Client
	SocialClient      socialityservice.Client
	InteractionClient interactionserver.Client
)
