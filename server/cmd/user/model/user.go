package model

import (
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type User struct {
	ID       int64  `gorm:"primarykey"`
	Username string `gorm:"type:varchar(33);unique;u;not null"`
	Password string `gorm:"type:varchar(33);not null"`
}

// BeforeCreate uses snowflake to generate an ID.
func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	sf, err := snowflake.NewNode(consts.UserSnowflakeNode)
	if err != nil {
		klog.Errorf("generate id failed: %s", err.Error())
		return err
	}
	u.ID = sf.Generate().Int64()
	return nil
}
