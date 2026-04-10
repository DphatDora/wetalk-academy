package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	_ "wetalk-academy/docs"
	"wetalk-academy/config"
	"wetalk-academy/internal/infrastructure/db"
	"wetalk-academy/internal/interface/router"
	"wetalk-academy/package/logger"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	DefaultPort = 8046
)

func main() {
	time.Local = time.UTC

	conf := config.GetConfig()
	if err := logger.Init(&conf.Log); err != nil {
		fmt.Fprintf(os.Stderr, "logger init: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = logger.Sync() }()

	logger.Debugf("[DEBUG] Config: %+v", conf)

	// init MongoDB
	mongoDB := db.NewMongoDB(&conf)
	defer func() {
		if err := mongoDB.Close(); err != nil {
			logger.Errorf("[ERROR] Failed to close MongoDB: %v", err)
		}
	}()

	// set up routes
	r := router.SetupRoutes(mongoDB.Database, &conf)

	// swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := conf.App.Port
	if port == 0 {
		port = DefaultPort
	}

	logger.Infof("[Server] Starting on PORT %d", port)
	r.Run(":" + strconv.Itoa(port))
}
