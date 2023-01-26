package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/model"
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
	createComment(db)
	createFavorite(db)
}

func createComment(db *gorm.DB) {
	_ = db.AutoMigrate(&model.Comment{})

	for i := 0; i < 10; i++ {
		c := model.Comment{
			UserId:      1616071000544256000,
			VideoId:     1616071000544256001,
			ActionType:  1,
			CommentText: fmt.Sprintf("It's comment%d", i),
			CreateDate:  time.Now(),
		}
		db.Save(&c)
	}
}

func createFavorite(db *gorm.DB) {
	_ = db.AutoMigrate(&model.Favorite{})

	for i := 0; i < 10; i++ {
		c := model.Favorite{
			UserId:     1616071000544256000,
			VideoId:    1616071000544256001,
			ActionType: 1,
		}
		db.Save(&c)
	}
}
