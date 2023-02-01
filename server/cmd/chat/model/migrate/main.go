package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/CyanAsterisk/TikGok/server/cmd/chat/model"
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

	_ = db.AutoMigrate(&model.Message{})

	for i := 0; i < 10; i++ {
		message := model.Message{
			ToUserId:   1616071000544256000,
			FromUserId: 1616071000577810432,
			Content:    fmt.Sprintf("contend%d", i),
			CreateDate: time.Now(),
		}
		db.Save(&message)
	}
}
