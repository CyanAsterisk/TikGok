package dao

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/model"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"gorm.io/gorm"
)

type Comment struct {
	db *gorm.DB
}

// NewComment create an interaction comment dao.
func NewComment(db *gorm.DB) *Comment {
	m := db.Migrator()
	if !m.HasTable(&model.Comment{}) {
		err := m.CreateTable(&model.Comment{})
		if err != nil {
			panic(err)
		}
	}
	return &Comment{
		db: db,
	}
}

// CommentCountByVideoId gets the number of comments by videoId.
func (c *Comment) CommentCountByVideoId(videoId int64) (int64, error) {
	var count int64
	err := c.db.Model(&model.Comment{}).
		Where(&model.Comment{VideoId: videoId, ActionType: consts.ValidComment}).Count(&count).Error
	if err != nil {
		return -1, err
	}
	return count, nil
}

// CommentIdListByVideoId gets commentId list by videoId
func (c *Comment) CommentIdListByVideoId(videoId int64) ([]string, error) {
	var commentIdList []string
	err := c.db.Model(&model.Comment{}).Select("id").
		Where(&model.Comment{VideoId: videoId}).Find(&commentIdList).Error
	if err != nil {
		return nil, err
	}
	return commentIdList, nil
}

// CreateComment creates a comment.
func (c *Comment) CreateComment(comment *model.Comment) (*model.Comment, error) {
	err := c.db.Model(model.Comment{}).
		Create(&comment).Error
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// DeleteComment to delete a comment.
func (c *Comment) DeleteComment(id int64) error {
	var comment model.Comment
	err := c.db.Model(model.Comment{}).
		Where(&model.Comment{ID: id, ActionType: consts.ValidComment}).First(&comment).Error
	if err != nil {
		return err
	}
	err = c.db.Model(model.Comment{}).
		Where(&model.Comment{ID: id}).Update("action_type", consts.InvalidComment).Error
	if err != nil {
		return err
	}
	return nil
}

// GetCommentListByVideoId gets comment list by videoId.
func (c *Comment) GetCommentListByVideoId(videoId int64) ([]*model.Comment, error) {
	var commentList []*model.Comment
	err := c.db.Model(model.Comment{}).
		Where(&model.Comment{VideoId: videoId, ActionType: consts.ValidComment}).Order("create_date desc").Find(&commentList).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return commentList, nil
}
