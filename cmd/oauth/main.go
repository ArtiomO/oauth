package main

import (

	"log"
	"github.com/ArtiomO/oauth/internal/server"

)

func main() {

	srv := &server.Server{}
	srv.InitGin()
	srv.InitCache()
	srv.InitClients()
	srv.SetupRouter()


	if err := srv.Gin.Run("0.0.0.0:8090"); err != nil {
		log.Printf("Error: %v", err)
	}
}
