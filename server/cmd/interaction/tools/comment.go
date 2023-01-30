package tools

import (
	"errors"
	"time"

	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/dao"
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/model"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction"
)

type CommentManager struct{}

// GetResp gets comment response.
func (m *CommentManager) GetResp(req *interaction.DouyinCommentActionRequest) (comment *model.Comment, err error) {
	switch req.ActionType {
	case consts.ValidComment:
		cmt, err := dao.CreateComment(&model.Comment{
			UserId:      req.UserId,
			VideoId:     req.VideoId,
			ActionType:  consts.ValidComment,
			CommentText: req.CommentText,
			CreateDate:  time.Now(),
		})
		if err != nil {
			return nil, err
		}
		return cmt, err
	case consts.InvalidComment:
		err = dao.DeleteComment(req.CommentId)
		if err != nil {
			return nil, err
		}
		return nil, nil
	default:
		return nil, errors.New("invalid action type")
	}
}
