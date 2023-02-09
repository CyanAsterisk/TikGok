package uploadService

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/disintegration/imaging"
	"github.com/minio/minio-go/v7"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Service struct {
	minioClient *minio.Client
	subscriber  *Subscriber
	publisher   *Publisher
}

func NewUploadService(minioClient *minio.Client, subscriber *Subscriber, publisher *Publisher) *Service {
	return &Service{
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

func (s *Service) UpLoadFile(videoFH *multipart.FileHeader) (playerUrl string, coverUrl string, err error) {
	suffix, err := getFileSuffix(videoFH.Filename)
	if err != nil {
		return "", "", err
	}
	if _, ok := videoSuffixSet[suffix]; !ok {
		return "", "", errno.ParamsEr.WithMessage("invalid video suffix")
	}

	sf, err := snowflake.NewNode(consts.MinioSnowflakeNode)
	if err != nil {
		hlog.Errorf("generate id failed: %s", err.Error())
		return "", "", err
	}
	taskId := sf.Generate().String()
	uploadPathBase := time.Now().Format("2006/01/02/") + taskId
	task := &VideoUploadTask{
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
	urlPrefix := consts.MinIOServer + "/" + consts.MinIOBucket + "/"
	return urlPrefix + task.VideoUploadPath, urlPrefix + task.CoverUploadPath, nil
}

func (s *Service) RunVideoUpload() error {
	taskCh, cleanUp, err := s.subscriber.Subscribe(context.Background())
	defer cleanUp()
	if err != nil {
		hlog.Error("cannot subscribe", err)
		return err
	}
	for task := range taskCh {
		if err = getVideoCover(task.VideoTmpPath, task.CoverTmpPath); err != nil {
			hlog.Errorf("get video cover err: videoTmpPath = %s", task.VideoTmpPath)
			continue
		}
		suffix, err := getFileSuffix(task.VideoTmpPath)
		if err != nil {
			hlog.Errorf("get video suffix err:videoTmpPath = %s", task.VideoTmpPath)
		}
		_, err = s.minioClient.FPutObject(context.Background(), consts.MinIOBucket, task.CoverUploadPath, task.CoverTmpPath, minio.PutObjectOptions{
			ContentType: "image/png",
		})
		_ = os.Remove(task.CoverTmpPath)
		_, err = s.minioClient.FPutObject(context.Background(), consts.MinIOBucket, task.VideoUploadPath, task.VideoTmpPath, minio.PutObjectOptions{
			ContentType: fmt.Sprintf("video/%s", suffix),
		})
		_ = os.Remove(task.VideoTmpPath)
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
