package main

import (
	"context"
	"github.com/ArtiomO/oauth/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	srv := &server.Server{}
	srv.InitGin()
	srv.InitCache()
	srv.InitClients()
	srv.SetupRouter()

	go func() {
		if err := srv.Gin.Run("0.0.0.0:8090"); err != nil {
			log.Printf("Error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	<-ctx.Done()
	log.Println("Server exiting")
}
