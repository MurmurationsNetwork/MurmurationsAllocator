package main

import (
	"context"
	"fmt"
	"github.com/MurmurationsNetwork/MurmurationsAllocator/config"
	"github.com/MurmurationsNetwork/MurmurationsAllocator/controllers"
	"github.com/MurmurationsNetwork/MurmurationsAllocator/database"
	"github.com/gin-contrib/cors"
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
	r.Use(cors.Default())
	
	r.GET("/", controllers.Ping)
	r.GET("/profiles", controllers.GetProfiles)
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Conf.Server.Port),
		Handler:      r,
		ReadTimeout:  config.Conf.Server.TimeoutRead,
		WriteTimeout: config.Conf.Server.TimeoutWrite,
		IdleTimeout:  config.Conf.Server.TimeoutIdle,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	database.DisconnectMongo()

	ctx, cancel := context.WithTimeout(context.Background(), config.Conf.Server.TimeoutIdle)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
