package upload_service

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/api/config"
	"github.com/CyanAsterisk/TikGok/server/cmd/api/global"
	"github.com/CyanAsterisk/TikGok/server/cmd/api/pkg/uploadService"
	"github.com/minio/minio-go/v7"
)

var (
	minioClient *minio.Client
	publisher   *uploadService.Publisher
	subscriber  *uploadService.Subscriber
	conf        *config.UploadServiceConfig
)

func Init() {
	initLogger()
	initConfig()
	initMinio()
	initMq()
	global.UploadService = uploadService.NewUploadService(minioClient, subscriber, publisher, conf)
}
