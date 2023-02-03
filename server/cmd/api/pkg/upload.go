package pkg

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/CyanAsterisk/TikGok/server/cmd/api/global"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/disintegration/imaging"
	"github.com/minio/minio-go/v7"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

var videoSuffixSet = map[string]struct{}{
	"mp4": {}, "avi": {}, "wmv": {},
	"w4v": {}, "asf": {}, "flv": {},
	"rmvb": {}, "rm": {}, "3gp": {},
	"vob": {}, "wma": {}, "mpeg": {},
	"mpg": {}, "mov": {},
}

func getVideoSuffix(header *multipart.FileHeader) (suffix string, err error) {
	fileName := header.Filename
	lastDotIndex := strings.LastIndex(fileName, ".")
	if lastDotIndex < 0 {
		return "", errors.New("missing suffix")
	}
	suffix = fileName[lastDotIndex+1:]
	suffix = strings.ToLower(suffix)
	if _, ok := videoSuffixSet[suffix]; !ok {
		return "", errors.New("suffix is invalid")
	}
	return suffix, nil
}

func UpLoadFile(videoFH *multipart.FileHeader) (playerUrl string, coverUrl string, err error) {
	suffix, err := getVideoSuffix(videoFH)
	if err != nil {
		return "", "", err
	}

	sf, err := snowflake.NewNode(consts.MinioSnowflakeNode)
	if err != nil {
		klog.Errorf("generate id failed: %s", err.Error())
		return "", "", err
	}
	vId := sf.Generate().String()
	videoTmpPath := "./tmp" + vId + "." + suffix
	coverTmpPath := "./tmp" + vId + ".png"

	videoFile, err := os.Create(videoTmpPath)
	if err != nil {
		return "", "", err
	}
	defer func() {
		videoFile.Close()
		os.Remove(videoTmpPath)
	}()

	mpFile, err := videoFH.Open()
	if err != nil {
		return "", "", err
	}
	defer mpFile.Close()

	_, err = videoFile.ReadFrom(mpFile)
	if err != nil {
		return "", "", err
	}

	if err = getVideoCover(videoTmpPath, coverTmpPath); err != nil {
		return "", "", err
	}
	defer os.Remove(coverTmpPath)

	uploadPathBase := time.Now().Format("2006/01/02") + vId
	videoUploadPath := uploadPathBase + "." + suffix
	coverUploadPath := uploadPathBase + ".png"

	_, err = global.MinioClient.FPutObject(context.Background(), consts.MinIOBucket, videoUploadPath, videoTmpPath, minio.PutObjectOptions{
		ContentType: fmt.Sprintf("video/%s", suffix),
	})
	if err != nil {
		return "", "", err
	}
	_, err = global.MinioClient.FPutObject(context.Background(), consts.MinIOBucket, coverUploadPath, coverTmpPath, minio.PutObjectOptions{
		ContentType: fmt.Sprintf("image/png"),
	})
	if err != nil {
		return "", "", err
	}
	urlPrefix := consts.MinIOServer + "/" + consts.MinIOBucket + "/"
	return urlPrefix + videoUploadPath, urlPrefix + coverUploadPath, nil
}

func getVideoCover(videoPath, coverPath string) error {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{"gte(n,1)"}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf).Run()
	if err != nil {
		return err
	}

	var img image.Image
	if img, err = imaging.Decode(buf); err != nil {
		return err
	}
	if err = imaging.Save(img, coverPath); err != nil {
		return err
	}
	return nil
}
