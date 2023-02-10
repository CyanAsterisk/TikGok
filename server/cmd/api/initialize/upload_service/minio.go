package upload_service

import (
	"context"

	"github.com/CyanAsterisk/TikGok/server/cmd/api/global"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func initMinio() {
	mi := global.ServerConfig.UploadServiceInfo.MinioInfo
	// Initialize minio client object.
	mc, err := minio.New(mi.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(mi.AccessKeyID, mi.SecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		hlog.Fatalf("create minio client err: %s", err.Error())
	}
	exists, err := mc.BucketExists(context.Background(), mi.Bucket)
	if err != nil {
		hlog.Fatal(err)
	}
	if !exists {
		err = mc.MakeBucket(context.Background(), mi.Bucket, minio.MakeBucketOptions{Region: "cn-north-1"})
		if err != nil {
			hlog.Fatalf("make bucket err: %s", err.Error())
		}
	}
	policy := `{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::` + mi.Bucket + `/*"],"Sid": ""}]}`
	err = mc.SetBucketPolicy(context.Background(), mi.Bucket, policy)
	if err != nil {
		hlog.Fatal("set bucket policy err:%s", err)
	}
	minioClient = mc
}
