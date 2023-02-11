package model

import (
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type Video struct {
	ID         int64  `gorm:"primarykey"`
	AuthorId   int64  `gorm:"column:author_id; not null"`
	PlayUrl    string `gorm:"not null; type: varchar(255)"`
	CoverUrl   string `gorm:"not null; type: varchar(255)"`
	Title      string `gorm:"not null; type: varchar(255)"`
	CreateTime int64  `gorm:"not null;"`
}

func (v *Video) BeforeCreate(_ *gorm.DB) (err error) {
	if v.ID == 0 {
		klog.Error("video id should be giving")
		return errno.VideoServerErr.WithMessage("video id should be giving")
	}
	return nil
}
