package app

import (
	"os"

	"github.com/PongDev/Go-Socket-Core/models"
	"github.com/PongDev/Go-Socket-Core/routes"
	"github.com/gin-gonic/gin"
)

type App struct {
	engine *gin.Engine
}

func NewApp() *App {
	trustProyy := os.Getenv("TRUST_PROXY")
	ginMode := os.Getenv("GIN_MODE")

	gin.SetMode(ginMode)
	app := &App{
		engine: gin.Default(),
	}
	app.engine.SetTrustedProxies([]string{trustProyy})
	return app
}

func (app *App) SetupRouter() {
	routes.SetupRouter(app.engine)
}

func (app *App) SetupDatabase() {
	models.ConnectDatabase()
}

func (app *App) Run() {
	app.engine.Run()
}
