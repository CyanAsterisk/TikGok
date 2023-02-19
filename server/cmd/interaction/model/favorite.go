package model

type Favorite struct {
	UserId     int64 `gorm:"not null"`
	VideoId    int64 `gorm:"not null"`
	ActionType int8  `gorm:"type:tinyint;not null"`
	CreateDate int64 `gorm:"not null"`
}
