package server

import (
	"context"
	"fmt"
	"os"

	"github.com/ArtiomO/oauth/internal/db"
	"github.com/ArtiomO/oauth/internal/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Server struct {
	DB      *pgx.Conn
	Gin     *gin.Engine
	Clients *[]db.Client
	Cache   repository.CacheRepository
}

func (s *Server) InitDb() *Server {

	db, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	s.DB = db
	return s
}

func (s *Server) InitCache() *Server {

	cache := &repository.RedisCacheRepository{}
	s.Cache = cache.InitRedisRepo()
	return s
}


func (s *Server) InitGin() *Server {

	r := gin.Default()
	r.LoadHTMLGlob("./web/templates/*")
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))
	r.Use(gin.Recovery())
	r.Static(os.Getenv("STATIC_URI"), "./web/static")
	s.Gin = r
	return s
}

func (s *Server) InitClients() *Server {

	clients := []db.Client{{
		ClientId:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURI:  "https://vertuhi.com/api/oauthcallback",
	}}
	s.Clients = &clients
	return s
}

func (s *Server) Ready() bool {
	return s.DB != nil && s.Gin != nil
}
