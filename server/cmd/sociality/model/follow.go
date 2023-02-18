package model

type Follow struct {
	UserId     int64 `gorm:"not null"`
	FollowerId int64 `gorm:"not null"`
	ActionType int8  `gorm:"type:tinyint;not null"`
}
