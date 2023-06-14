package server

import (
	_ "github.com/ArtiomO/oauth/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) SetupRouter() {

	s.Gin.GET("/api/v1/login", s.GetLoginHandler)
	s.Gin.POST("/api/v1/login", s.PostLoginHandler)
	s.Gin.POST("/api/v1/token", s.PostTokenHandler)
	s.Gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

}
