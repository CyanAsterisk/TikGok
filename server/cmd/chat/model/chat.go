package model

import (
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type Message struct {
	ID         int64 `gorm:"primarykey"`
	UserId     int64 `gorm:"not null"`
	FollowerId int64 `gorm:"not null"`
	ActionType int8  `gorm:"type:tinyint;not null"`
}

// BeforeCreate uses snowflake to generate an ID.
func (m *Message) BeforeCreate(_ *gorm.DB) (err error) {
	sf, err := snowflake.NewNode(consts.ChatSnowflakeNode)
	if err != nil {
		klog.Errorf("generate id failed: %s", err.Error())
		return err
	}
	m.ID = sf.Generate().Int64()
	return nil
}
