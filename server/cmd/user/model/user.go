package model

type User struct {
	ID              int64  `gorm:"primarykey"`
	Username        string `gorm:"type:varchar(33);unique;not null"`
	Password        string `gorm:"type:varchar(33);not null"`
	Avatar          string `gorm:"type:varchar(255);not null"`
	BackGroundImage string `gorm:"type:varchar(255);not null"`
	Signature       string `gorm:"type:varchar(255);not null"`
}
