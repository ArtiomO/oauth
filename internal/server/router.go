package server

func (s *Server) SetupRouter() {

	s.Gin.GET("/api/v1/login", s.GetLoginHandler)
	s.Gin.POST("/api/v1/login", s.PostLoginHandler)
	s.Gin.POST("/api/v1/token", s.PostTokenHandler)

}
