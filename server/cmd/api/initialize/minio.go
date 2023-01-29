package initialize

import (
	"context"

	"github.com/CyanAsterisk/TikGok/server/cmd/api/global"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/u2takey/go-utils/klog"
)

func InitMinio() {
	mi := global.ServerConfig.MinioInfo
	// Initialize minio client object.
	mc, err := minio.New(mi.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(mi.AccessKeyID, mi.SecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		klog.Fatalln(err)
	}
	bucketName := consts.MinIOBucket
	exists, err := mc.BucketExists(context.Background(), bucketName)
	if err != nil {
		klog.Fatalln(err)
	}
	if !exists {
		err = mc.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: "cn-north-1"})
		if err != nil {
			klog.Fatalln(err)
		}
	}
	policy := `{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::` + bucketName + `/*"],"Sid": ""}]}`
	err = mc.SetBucketPolicy(context.Background(), bucketName, policy)
	if err != nil {
		klog.Fatalln(err)
	}
	global.MinioClient = mc
}
