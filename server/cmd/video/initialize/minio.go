package initialize

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/video/global"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/minio/minio-go"
	"github.com/minio/minio-go/v6/pkg/policy"
)

func InitMinio() {
	mi := global.ServerConfig.MinioInfo
	// 初使化 minio client对象。
	mc, err := minio.New(mi.Endpoint, mi.AccessKeyID, mi.SecretAccessKey, false)
	if err != nil {
		klog.Fatalf("create minioClient err:%s", err.Error())
	}
	err = mc.SetBucketPolicy(consts.MinIOBucket, policy.BucketPolicyReadWrite)
	if err != nil {
		klog.Fatalf("set bucket policy err: %s", err.Error())
	}
	global.MinioClient = mc
}
