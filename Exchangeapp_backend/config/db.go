package config

import (
	"exchangeapp/global"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDB(){
	dsn := AppConfig.Database.Dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
  
	if err!=nil{
		log.Printf("Failed to initialize database, got error: %v", err)
		return
	}

	sqlDB, err := db.DB()

	if err !=nil{
		log.Printf("Failed to configure database, got error: %v", err)
		return
	}

	sqlDB.SetMaxIdleConns(AppConfig.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(AppConfig.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	global.Db = db
	log.Println("Database initialized successfully")
}