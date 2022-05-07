package main

import (
	"fmt"
	"github.com/ashishkumar68/auction-api/config"
	"github.com/ashishkumar68/auction-api/middleware"
	"github.com/ashishkumar68/auction-api/migrations"
	"github.com/ashishkumar68/auction-api/routes"
	"github.com/ashishkumar68/auction-api/validators"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func SetupRoutes(engine *gin.Engine) {
	apiGroup := engine.Group("/api", middleware.CorsRoute())
	routes.MapIndexRoutes(engine)
	routes.MapUserRoutes(apiGroup)
	routes.MapItemRoutes(apiGroup)
}

func runMigrations() {
	migrations.DropAndCreateTables()
}

func main() {
	config.InitialiseConfig()
	// creating gin engine.
	engine := gin.Default()
	_ = engine.SetTrustedProxies(nil)
	runMigrations()
	// setting up routes.
	SetupRoutes(engine)
	// setting up custom validators.
	validators.SetupCustomValidators()

	// launch server.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	err := engine.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalln("could not start server", err)
	}
	log.Println(fmt.Sprintf("server listening on port: %s", port))
}
