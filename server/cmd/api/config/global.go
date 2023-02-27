package config

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/api/pkg/uploadService"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/chat/chatservice"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction/interactionserver"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/sociality/socialityservice"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/user/userservice"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/video/videoservice"
)

var (
	GlobalServerConfig  = &ServerConfig{}
	GlobalNacosConfig   = &NacosConfig{}
	GlobalUploadService *uploadService.Service

	GlobalChatClient        chatservice.Client
	GlobalUserClient        userservice.Client
	GlobalVideoClient       videoservice.Client
	GlobalSocialClient      socialityservice.Client
	GlobalInteractionClient interactionserver.Client
)
