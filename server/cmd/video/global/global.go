package global

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/video/config"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction/interactionserver"

	"github.com/minio/minio-go"
	"gorm.io/gorm"
)

var (
	DB             *gorm.DB
	ServerConfig   config.ServerConfig
	NacosConfig    config.NacosConfig
	MinioClient    *minio.Client
	InteractClient interactionserver.Client
)
