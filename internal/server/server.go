package server

import (
	"os"

	"github.com/ArtiomO/oauth/internal/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Gin     *gin.Engine
	Clients repository.ClientsRepository
	Cache   repository.CacheRepository
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

	clientsRepo := &repository.MemoryClientRepository{}
	s.Clients = clientsRepo.InitClientRepo()
	return s
}
