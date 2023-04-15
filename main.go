package main

import (
	"github.com/PongDev/Go-Socket-Core/app"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load .env file")
	}

	app := app.NewApp()
	app.SetupRouter()
	app.SetupDatabase()
	app.Run()
}
