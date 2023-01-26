package pack

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/model"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction"
)

// Comment model to idl
func Comment(c *model.Comment) *interaction.Comment {
	if c == nil {
		return nil
	}
	return &interaction.Comment{
		Id: c.ID,
		User: &interaction.User{
			Id: c.ID,
		},
		Content:    c.CommentText,
		CreateDate: c.CreateDate.Format("mm-dd"),
	}
}

// Comments model to idl
func Comments(c []*model.Comment) []*interaction.Comment {
	if c == nil {
		return nil
	}
	cl := make([]*interaction.Comment, 0)
	for _, cmt := range c {
		cl = append(cl, Comment(cmt))
	}
	return cl
}
