package pkg

import "github.com/CyanAsterisk/TikGok/server/cmd/interaction/model"

type FavoriteRedisManager struct {
}

func (r *FavoriteRedisManager) FavoriteCountByVideoId(videoId int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}
func (r *FavoriteRedisManager) CreateFavorite(fav *model.Favorite) error {
	//TODO implement me
	panic("implement me")
}
func (r *FavoriteRedisManager) UpdateFavorite(userId, videoId int64, actionType int8) error {
	//TODO implement me
	panic("implement me")
}
func (r *FavoriteRedisManager) GetFavoriteInfo(userId, videoId int64) (*model.Favorite, error) {
	//TODO implement me
	panic("implement me")
}
func (r *FavoriteRedisManager) GetFavoriteVideoIdListByUserId(userId int64) ([]int64, error) {
	//TODO implement me
	panic("implement me")
}
