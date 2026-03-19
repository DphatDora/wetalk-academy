package config

import (
	"log"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var (
	once   sync.Once
	config Config
)

type Config struct {
	App      App
	Database Database
	Log      Log
	Server   Server
	Auth     Auth
	Client   Client
	Judge0   Judge0
}

func LoadConfig() {
	once.Do(func() {
		// load .env config
		_ = godotenv.Load()

		// load yaml config
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./config")

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Config file error: %s", err)
		}

		// bind system environment variables
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		bindEnvs()

		// load into struct
		if err := viper.Unmarshal(&config); err != nil {
			log.Fatalf("Config unmarshal error: %s", err)
		}
	})
}

func GetConfig() Config {
	LoadConfig()
	return config
}

func bindEnvs() {
	// App
	_ = viper.BindEnv("app.name", "APP_NAME")
	_ = viper.BindEnv("app.host", "HOST")
	_ = viper.BindEnv("app.port", "PORT")

	// Database
	_ = viper.BindEnv("database.uri", "MONGO_URI")
	_ = viper.BindEnv("database.name", "MONGO_DB_NAME")

	// Server
	_ = viper.BindEnv("server.url", "SERVER_URL")

	// Client
	_ = viper.BindEnv("client.url", "CLIENT_URL")

	// Judge0
	_ = viper.BindEnv("judge0.baseURL", "JUDGE0_BASE_URL")

	// Auth
	_ = viper.BindEnv("auth.jwtSecret", "JWT_SECRET")
}
