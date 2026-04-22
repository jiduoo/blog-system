package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"backend/config"
	"backend/global"
	"backend/models"
	"backend/router"
	"backend/utils"
)

func main() {
	// 初始化配置（数据库、Redis等）
	config.InitConfig()

	// 自动迁移所有模型到数据库
	migrateModels()

	// 初始化root用户（如果不存在）
	initRootUser()

	// 配置并启动HTTP服务器
	r := router.SetupRouter()

	port := config.AppConfig.App.Port
	if port == "" {
		port = ":8080"
	}

	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	// 启动服务器在后台goroutine中运行
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	// 给服务器5秒时间完成正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

// migrateModels 自动迁移所有数据模型到数据库
// GORM会自动创建表和更新表结构
func migrateModels() {
	if err := global.Db.AutoMigrate(
		&models.User{},
		&models.Blog{},
		&models.Tag{},
	); err != nil {
		log.Printf("Failed to migrate models: %v", err)
		return
	}
	log.Println("Models migrated successfully")
}

// initRootUser 初始化root用户
// 如果数据库中没有root用户，则创建一个root用户用于管理
func initRootUser() {
	var rootUser models.User
	result := global.Db.Where("username = ?", "root").First(&rootUser)

	if result.Error != nil {
		hashedPwd, err := utils.HashPassword("root")
		if err != nil {
			log.Printf("Failed to hash root password: %v", err)
			return
		}

		newRootUser := models.User{
			Username: "root",
			Password: hashedPwd,
			IsRoot:   true,
		}

		if err := global.Db.Create(&newRootUser).Error; err != nil {
			log.Printf("Failed to create root user: %v", err)
			return
		}

		log.Println("Root user created: username=root, password=root")
	}
}
