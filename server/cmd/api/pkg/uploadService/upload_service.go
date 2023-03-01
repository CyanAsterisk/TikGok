package uploadService

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"mime/multipart"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/CyanAsterisk/TikGok/server/cmd/api/pkg/uploadService/config"
	"github.com/CyanAsterisk/TikGok/server/cmd/api/pkg/uploadService/pkg"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/disintegration/imaging"
	"github.com/minio/minio-go/v7"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Service struct {
	config      *config.UploadServiceConfig
	minioClient *minio.Client
	subscriber  *pkg.Subscriber
	publisher   *pkg.Publisher
}

func NewUploadService(minioClient *minio.Client, subscriber *pkg.Subscriber, publisher *pkg.Publisher, config *config.UploadServiceConfig) *Service {
	return &Service{
		config:      config,
		minioClient: minioClient,
		subscriber:  subscriber,
		publisher:   publisher,
	}
}

var videoSuffixSet = map[string]struct{}{
	"mp4": {}, "avi": {}, "wmv": {},
	"w4v": {}, "asf": {}, "flv": {},
	"rmvb": {}, "rm": {}, "3gp": {},
	"vob": {}, "wma": {}, "mpeg": {},
	"mpg": {}, "mov": {},
}

func (s *Service) UpLoadFile(videoFH *multipart.FileHeader) (playerUrl, coverUrl string, err error) {
	suffix, err := getFileSuffix(videoFH.Filename)
	if err != nil {
		return "", "", err
	}
	if _, ok := videoSuffixSet[suffix]; !ok {
		return "", "", errno.ParamsEr.WithMessage("invalid video suffix")
	}

	sf, err := snowflake.NewNode(consts.MinioSnowflakeNode)
	if err != nil {
		klog.Errorf("generate id failed: %s", err.Error())
		return "", "", err
	}
	taskId := sf.Generate().String()
	uploadPathBase := time.Now().Format("2006/01/02/") + taskId
	task := &pkg.VideoUploadTask{
		VideoTmpPath:    "./tmp" + taskId + "." + suffix,
		CoverTmpPath:    "./tmp" + taskId + ".png",
		VideoUploadPath: uploadPathBase + "." + suffix,
		CoverUploadPath: uploadPathBase + ".png",
	}
	videoFile, err := os.Create(task.VideoTmpPath)
	if err != nil {
		return "", "", err
	}
	defer videoFile.Close()

	mpFile, err := videoFH.Open()
	if err != nil {
		return "", "", err
	}
	defer mpFile.Close()

	_, err = videoFile.ReadFrom(mpFile)
	if err != nil {
		return "", "", err
	}
	if err = s.publisher.Publish(context.Background(), task); err != nil {
		return "", "", err
	}
	urlPrefix := s.config.MinioInfo.UrlPrefix
	return urlPrefix + task.VideoUploadPath, urlPrefix + task.CoverUploadPath, nil
}

func (s *Service) RunVideoUpload() error {
	taskCh, cleanUp, err := s.subscriber.Subscribe(context.Background())
	defer cleanUp()
	if err != nil {
		klog.Error("cannot subscribe", err)
		return err
	}
	for task := range taskCh {
		if err = getVideoCover(task.VideoTmpPath, task.CoverTmpPath); err != nil {
			klog.Errorf("get video cover err: videoTmpPath = %s", task.VideoTmpPath)
			continue
		}
		buckName := s.config.MinioInfo.Bucket

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer func() {
				wg.Done()
				_ = os.Remove(task.CoverTmpPath)
			}()
			if _, err = s.minioClient.FPutObject(context.Background(), buckName, task.CoverUploadPath, task.CoverTmpPath, minio.PutObjectOptions{
				ContentType: "image/png",
			}); err != nil {
				klog.Error("upload cover image err", err)
			}
		}()

		go func() {
			defer func() {
				wg.Done()
				_ = os.Remove(task.VideoTmpPath)
			}()
			suffix, err := getFileSuffix(task.VideoTmpPath)
			if err != nil {
				klog.Errorf("get video suffix err:videoTmpPath = %s", task.VideoTmpPath)
				return
			}
			if _, err = s.minioClient.FPutObject(context.Background(), buckName, task.VideoUploadPath, task.VideoTmpPath, minio.PutObjectOptions{
				ContentType: fmt.Sprintf("video/%s", suffix),
			}); err != nil {
				klog.Error("upload video err", err)
			}
		}()
		wg.Wait()
	}
	return nil
}

func getFileSuffix(fileName string) (suffix string, err error) {
	lastDotIndex := strings.LastIndex(fileName, ".")
	if lastDotIndex < 0 {
		return "", errno.ParamsEr.WithMessage("missing suffix")
	}
	suffix = fileName[lastDotIndex+1:]
	suffix = strings.ToLower(suffix)
	return suffix, nil
}

// save video cover image to giving path.
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
