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
	setUpInfrastructure()
	defer closeInfrastructure()
}

func setUpInfrastructure() {
	time.Local = time.UTC

	conf := config.GetConfig()
	fmt.Println("[DEBUG] Config:", conf)

	// init MongoDB
	db.InitMongoDB(&conf)

	// set up routes
	r := router.SetupRoutes(db.GetDB(), &conf)

	port := conf.App.Port
	if port == 0 {
		port = DefaultPort
	}

	log.Printf("[Server] Starting on PORT %d", port)
	r.Run(":" + strconv.Itoa(port))
}

func closeInfrastructure() {
	if err := db.CloseMongoDB(); err != nil {
		log.Printf("[ERROR] Close MongoDB fail: %s\n", err)
	}
}
