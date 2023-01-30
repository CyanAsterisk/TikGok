package global

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/api/config"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction/interactionserver"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/sociality/socialityservice"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/user/userservice"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/video/videoservice"
	"github.com/minio/minio-go/v7"
)

var (
	ServerConfig = &config.ServerConfig{}
	NacosConfig  = &config.NacosConfig{}

	MinioClient       *minio.Client
	UserClient        userservice.Client
	VideoClient       videoservice.Client
	SocialClient      socialityservice.Client
	InteractionClient interactionserver.Client
)
