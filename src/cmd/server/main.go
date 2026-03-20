package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"wetalk-academy/config"
	"wetalk-academy/internal/infrastructure/db"
	"wetalk-academy/internal/interface/router"
)

const (
	DefaultPort = 8046
)

func main() {
	time.Local = time.UTC

	conf := config.GetConfig()
	fmt.Println("[DEBUG] Config:", conf)

	// init MongoDB
	mongoDB := db.NewMongoDB(&conf)
	defer func() {
		if err := mongoDB.Close(); err != nil {
			log.Printf("[ERROR] Failed to close MongoDB: %v", err)
		}
	}()

	// set up routes
	r := router.SetupRoutes(mongoDB.Database, &conf)

	port := conf.App.Port
	if port == 0 {
		port = DefaultPort
	}

	log.Printf("[Server] Starting on PORT %d", port)
	r.Run(":" + strconv.Itoa(port))
}
