package main

import (
	"fmt"
	"github.com/ashishkumar68/auction-api/config"
	"github.com/ashishkumar68/auction-api/middleware"
	"github.com/ashishkumar68/auction-api/migrations"
	"github.com/ashishkumar68/auction-api/routes"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func SetupRoutes(engine *gin.Engine) {
	apiGroup := engine.Group(
		"/api", middleware.AuthenticatedRoute(), middleware.TransactionRoute(),
	)
	routes.MapIndexRoutes(engine)
	routes.MapAuthRoutes(apiGroup)
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
