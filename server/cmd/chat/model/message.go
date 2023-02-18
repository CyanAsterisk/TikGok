package model

import (
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type Message struct {
	ID         int64  `gorm:"primarykey"`
	ToUserId   int64  `gorm:"not null"`
	FromUserId int64  `gorm:"not null"`
	Content    string `gorm:"type:varchar(256);not null"`
	CreateTime int64  `gorm:"not null"`
}

// BeforeCreate uses snowflake to generate an ID.
func (m *Message) BeforeCreate(_ *gorm.DB) (err error) {
	// skip if id already set
	if m.ID != 0 {
		return nil
	}
	sf, err := snowflake.NewNode(consts.ChatSnowflakeNode)
	if err != nil {
		klog.Errorf("generate id failed: %s", err.Error())
		return err
	}
	m.ID = sf.Generate().Int64()
	return nil
}
