package main

import (
	"context"
	"fmt"
	"github.com/ashishkumar68/auction-api/config"
	"github.com/ashishkumar68/auction-api/migrations"
	"github.com/ashishkumar68/auction-api/routes"
	"github.com/ashishkumar68/auction-api/validators"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func SetupRoutes(engine *gin.Engine) {
	apiGroup := engine.Group("/api")

	routes.MapIndexRoutes(engine)
	routes.MapUserRoutes(apiGroup)
	routes.MapItemRoutes(apiGroup)
}

func runMigrations() {
	migrations.RunMigrations()
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

	if config.AppEnvTest == os.Getenv("APP_ENV") {
		launchDevelopmentServer(engine)
	} else {
		launchReleaseServer(engine)
	}
}

func launchDevelopmentServer(engine *gin.Engine) {
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

func launchReleaseServer(engine *gin.Engine) {
	// launch server.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	port = fmt.Sprintf(":%s", port)
	server := &http.Server{
		Addr:    port,
		Handler: engine,
	}

	go func() {
		log.Println(fmt.Sprintf("server listening on port: %s", port))
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalln("could not start server", err)
		}
	}()

	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Initiating Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Error while Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server gracefully shutdown.")
}
