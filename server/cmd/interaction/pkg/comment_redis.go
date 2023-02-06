package pkg

import "github.com/CyanAsterisk/TikGok/server/cmd/interaction/model"

type CommentRedisManager struct {
}

func (r *CommentRedisManager) CommentCountByVideoId(videoId int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}
func (r *CommentRedisManager) CreateComment(comment *model.Comment) (*model.Comment, error) {
	//TODO implement me
	panic("implement me")
}
func (r *CommentRedisManager) DeleteComment(id int64) error {
	//TODO implement me
	panic("implement me")
}
func (r *CommentRedisManager) GetCommentListByVideoId(videoId int64) ([]*model.Comment, error) {
	//TODO implement me
	panic("implement me")
}
