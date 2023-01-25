package model

import (
	"time"

	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type Comment struct {
	ID          int64     `gorm:"primarykey"`
	UserId      int64     `gorm:"not null"`
	VideoId     int64     `gorm:"not null"`
	ActionType  int8      `gorm:"type:tinyint;not null"`
	CommentText string    `gorm:"type:varchar(256);not null"`
	CreateDate  time.Time `gorm:"not null"`
}

// BeforeCreate uses snowflake to generate an ID.
func (c *Comment) BeforeCreate(_ *gorm.DB) (err error) {
	sf, err := snowflake.NewNode(consts.InteractionSnowflakeNode)
	if err != nil {
		klog.Errorf("generate id failed: %s", err.Error())
		return err
	}
	c.ID = sf.Generate().Int64()
	return nil
}
