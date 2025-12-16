package main

import (
	"fmt"

	"github.com/fsangopanta/demo-soft-code/config"
	"github.com/fsangopanta/demo-soft-code/middlewares"
	service_desk "github.com/fsangopanta/demo-soft-code/modules/service_desk"
	providers "github.com/fsangopanta/demo-soft-code/providers"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

func run(server *gin.Engine) {
	cfg := config.Load()
	port := cfg.App.Port

	var serve string
	if cfg.App.Environment == "dev" {
		serve = "0.0.0.0:" + fmt.Sprint(port)
	} else {
		serve = ":" + fmt.Sprint(port)
	}

	if cfg.App.Environment == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	
	if err := server.Run(serve); err != nil {
		fmt.Printf("error running server: %v", err)
	}
}



func main() {
	var (
		injector = do.New()
	)

	providers.RegisterDependencies(injector)

	// TODO: Manejar inyeccion de dependencias
	server := gin.Default()
	server.Use(middlewares.CORSMiddleware())
	// Register incoming modules
	service_desk.RegisterRoutes(server, injector)

	
	run(server)
}