package main

import (
	"github.com/greenac/sqsmock/app"
	"github.com/greenac/sqsmock/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
	}

	app.Start()
}
