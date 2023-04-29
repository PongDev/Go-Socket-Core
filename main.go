package main

import (
	"log"

	"github.com/PongDev/Go-Socket-Core/app"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed to load .env file, Use from environment variables")
	}

	app := app.NewApp()
	app.SetupRouter()
	app.SetupDatabase()
	app.Run()
}
