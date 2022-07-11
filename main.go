package main

import (
	"log"
	"os"
	"test/test_app/app/api/server"
	"test/test_app/app/service/logger"
	"test/test_app/config"

	"github.com/joho/godotenv"
)

func init() {
	environmentPath := ".env"
	err := godotenv.Load(environmentPath)
	if err != nil {
		log.Fatal("Failed to Load env", err)
	}
}
func main() {

	serviceName := "test"
	environment := os.Getenv("BOOT_CUR_ENV")
	if environment == "" {
		environment = "dev"
	}

	config.Init(serviceName, environment)
	logger.InitLogger()

	server.Init()
}
