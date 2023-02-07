package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/CyanAsterisk/TikGok/server/cmd/video/model"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func main() {
	// Defined by your database.
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL Threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color printing
		},
	)

	// global mode
	db, err := gorm.Open(mysql.Open(consts.UserMigrateDSN), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	_ = db.AutoMigrate(&model.Video{})

	for i := 0; i < 10; i++ {
		video := model.Video{
			ID:         0,
			AuthorId:   int64(10000000000 + i),
			PlayUrl:    "fake-playUrl",
			CoverUrl:   "fake-playUrl",
			Title:      fmt.Sprintf("fake-title-%d", i),
			CreateTime: time.Now().UnixNano(),
		}
		db.Save(&video)
	}
}
