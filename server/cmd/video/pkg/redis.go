package pkg

import "github.com/CyanAsterisk/TikGok/server/cmd/video/model"

type RedisManager struct {
}

func (r *RedisManager) CreateVideo(video *model.Video) error {
	//TODO implement me
	panic("implement me")
}
func (r *RedisManager) GetVideosByLatestTime(latestTime int64) ([]*model.Video, error) {
	//TODO implement me
	panic("implement me")
}
func (r *RedisManager) GetVideosByUserId(uid int64) ([]*model.Video, error) {
	//TODO implement me
	panic("implement me")
}
func (r *RedisManager) GetVideoByVideoId(vid int64) (*model.Video, error) {
	//TODO implement me
	panic("implement me")
}
func (r *RedisManager) BatchGetVideoByVideoId(vidList []int64) ([]*model.Video, error) {
	//TODO implement me
	panic("implement me")
}
