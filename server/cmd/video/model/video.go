package model

import (
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	Uid      int64  `gorm:"column:user_id; not null"`
	PlayUrl  string `gorm:"not null; type: varchar(255)"`
	CoverUrl string `gorm:"not null; type: varchar(255)"`
	Title    string `gorm:"not null; type: varchar(255)"`
}

// BeforeCreate uses snowflake to generate an ID.
func (v *Video) BeforeCreate(_ *gorm.DB) (err error) {
	sf, err := snowflake.NewNode(consts.VideoSnowflakeNode)
	if err != nil {
		klog.Errorf("generate id failed: %s", err.Error())
		return err
	}
	v.ID = uint(sf.Generate().Int64())
	return nil
}
