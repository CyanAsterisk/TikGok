package tools

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/model"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
)

// Comment model to idl
func Comment(c *model.Comment) *base.Comment {
	if c == nil {
		return nil
	}
	return &base.Comment{
		Id: c.ID,
		User: &base.User{
			Id: c.ID,
		},
		Content:    c.CommentText,
		CreateDate: c.CreateDate.Format("mm-dd"),
	}
}

// Comments model to idl
func Comments(c []*model.Comment) []*base.Comment {
	if c == nil {
		return nil
	}
	cl := make([]*base.Comment, 0)
	for _, cmt := range c {
		cl = append(cl, Comment(cmt))
	}
	return cl
}
