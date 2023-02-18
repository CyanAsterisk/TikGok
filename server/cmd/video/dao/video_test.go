package dao

import (
	"fmt"
	"testing"

	"github.com/CyanAsterisk/TikGok/server/cmd/video/model"
	"github.com/CyanAsterisk/TikGok/server/shared/test"
	"github.com/bytedance/sonic"
)

func TestVideoLifecycle(t *testing.T) {

	cleanUpFunc, db, err := test.RunWithMySQLInDocker(t)
	defer cleanUpFunc()
	if err != nil {
		t.Fatal(err)
	}
	dao := NewVideo(db)

	timeStamp := int64(1676323214)
	videoList := make([]*model.Video, 0)
	for i := int64(0); i < 5; i++ {
		v := &model.Video{
			ID:         100000 + i,
			AuthorId:   200000 + i%2,
			PlayUrl:    fmt.Sprintf("vidoe%d-fake-play-url", i),
			CoverUrl:   fmt.Sprintf("vidoe%d-fake-cover-url", i),
			Title:      fmt.Sprintf("video%d-tiltle", i),
			CreateTime: timeStamp + i,
		}
		videoList = append(videoList, v)
	}
	cases := []struct {
		name       string
		op         func() (string, error)
		wantErr    bool
		wantResult string
	}{
		{
			name: "create video",
			op: func() (string, error) {
				for _, v := range videoList {
					if err = dao.CreateVideo(v); err != nil {
						return "", err
					}
				}
				return "", nil
			},
			wantErr:    false,
			wantResult: "",
		},
		{
			name: "duplicate create video",
			op: func() (string, error) {
				err = dao.CreateVideo(videoList[0])
				return "", err
			},
			wantErr: true,
		},
		{
			name: "get video by id",
			op: func() (string, error) {
				video, err := dao.GetVideoByVideoId(videoList[0].ID)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(video)
				if err != nil {
					return "", err
				}
				return string(result), nil
			},
			wantErr:    false,
			wantResult: `{"ID":100000,"AuthorId":200000,"PlayUrl":"vidoe0-fake-play-url","CoverUrl":"vidoe0-fake-cover-url","Title":"video0-tiltle","CreateTime":1676323214}`,
		},
		{
			name: "get videoList by Author id",
			op: func() (string, error) {
				video, err := dao.GetVideoListByAuthorId(videoList[0].AuthorId)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(video)
				if err != nil {
					return "", err
				}
				return string(result), nil
			},
			wantErr:    false,
			wantResult: `[{"ID":100004,"AuthorId":200000,"PlayUrl":"vidoe4-fake-play-url","CoverUrl":"vidoe4-fake-cover-url","Title":"video4-tiltle","CreateTime":1676323218},{"ID":100002,"AuthorId":200000,"PlayUrl":"vidoe2-fake-play-url","CoverUrl":"vidoe2-fake-cover-url","Title":"video2-tiltle","CreateTime":1676323216},{"ID":100000,"AuthorId":200000,"PlayUrl":"vidoe0-fake-play-url","CoverUrl":"vidoe0-fake-cover-url","Title":"video0-tiltle","CreateTime":1676323214}]`,
		},
		{
			name: "get videoIdList by Author id",
			op: func() (string, error) {
				video, err := dao.GetVideoIdListByAuthorId(videoList[0].AuthorId)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(video)
				if err != nil {
					return "", err
				}
				return string(result), nil
			},
			wantErr:    false,
			wantResult: `[100000,100002,100004]`,
		},
		{
			name: "get video list by last time",
			op: func() (string, error) {
				videos, err := dao.GetVideoListByLatestTime(timeStamp + 2)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(videos)
				if err != nil {
					return "", err
				}
				return string(result), err
			},
			wantErr:    false,
			wantResult: `[{"ID":100002,"AuthorId":200000,"PlayUrl":"vidoe2-fake-play-url","CoverUrl":"vidoe2-fake-cover-url","Title":"video2-tiltle","CreateTime":1676323216},{"ID":100001,"AuthorId":200001,"PlayUrl":"vidoe1-fake-play-url","CoverUrl":"vidoe1-fake-cover-url","Title":"video1-tiltle","CreateTime":1676323215},{"ID":100000,"AuthorId":200000,"PlayUrl":"vidoe0-fake-play-url","CoverUrl":"vidoe0-fake-cover-url","Title":"video0-tiltle","CreateTime":1676323214}]`,
		},
		{
			name: "batch get video by id",
			op: func() (string, error) {
				videoList, err := dao.BatchGetVideoByVideoId([]int64{videoList[0].ID, videoList[1].ID})
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(videoList)
				if err != nil {
					return "", err
				}
				return string(result), nil
			},
			wantErr:    false,
			wantResult: `[{"ID":100000,"AuthorId":200000,"PlayUrl":"vidoe0-fake-play-url","CoverUrl":"vidoe0-fake-cover-url","Title":"video0-tiltle","CreateTime":1676323214},{"ID":100001,"AuthorId":200001,"PlayUrl":"vidoe1-fake-play-url","CoverUrl":"vidoe1-fake-cover-url","Title":"video1-tiltle","CreateTime":1676323215}]`,
		},
		{
			name: "delete video by id",
			op: func() (string, error) {
				err := dao.DeleteVideoById(videoList[0].ID)
				return "", err
			},
			wantErr:    false,
			wantResult: "",
		},
		{
			name: "duplicate delete video by id",
			op: func() (string, error) {
				err := dao.DeleteVideoById(videoList[0].ID)
				return "", err
			},
			wantErr: true,
		},
	}

	for _, cc := range cases {
		result, err := cc.op()
		if cc.wantErr {
			if err == nil {
				t.Errorf("%s:want error;got none", cc.name)
			} else {
				continue
			}
		}
		if err != nil {
			t.Errorf("%s:operation failed: %v", cc.name, err)
		}
		if result != cc.wantResult {
			t.Errorf("%s:result err: want %s,got %s", cc.name, cc.wantResult, result)
		}
	}
}
