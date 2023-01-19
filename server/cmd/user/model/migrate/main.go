package main

import (
	"fmt"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"log"
	"os"
	"time"

	"github.com/CyanAsterisk/TikGok/server/cmd/user/model"
	"github.com/CyanAsterisk/TikGok/server/cmd/user/tools"
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

	_ = db.AutoMigrate(&model.User{})

	for i := 0; i < 10; i++ {
		cryPassword := tools.Md5Crypt(fmt.Sprintf("Password%d", i), "TikGok")
		user := model.User{
			Username: fmt.Sprintf("User%d", i),
			Password: cryPassword,
		}
		db.Save(&user)
	}
}
