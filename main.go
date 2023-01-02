package main

import (
	"MurmurationsAllocator/config"
	"MurmurationsAllocator/controllers"
	"MurmurationsAllocator/database"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	config.Init()
	database.ConnectMongo()
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/", controllers.Ping)
	r.GET("/profiles", controllers.GetProfiles)
	svc := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Conf.Server.Port),
		Handler:      r,
		ReadTimeout:  config.Conf.Server.TimeoutRead,
		WriteTimeout: config.Conf.Server.TimeoutWrite,
		IdleTimeout:  config.Conf.Server.TimeoutIdle,
	}

	closed := make(chan struct{})
	go waitForShutdown(svc, closed)

	<-closed
}

func waitForShutdown(server *http.Server, closed chan struct{}) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	database.DisconnectMongo()

	ctx, cancel := context.WithTimeout(context.Background(), config.Conf.Server.TimeoutIdle)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Index service shutdown failure", err)
	}

	close(closed)
}
