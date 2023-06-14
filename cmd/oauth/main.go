package main

import (
	"github.com/ArtiomO/oauth/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	defer srv.Cache.Disconnect()

}
