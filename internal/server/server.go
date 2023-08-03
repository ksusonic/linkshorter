package server

import "github.com/gin-gonic/gin"

type Server struct {
	*gin.Engine
}

func NewServer() *Server {
	engine := gin.Default()
	_ = engine.SetTrustedProxies([]string{})
	return &Server{Engine: engine}
}

func (s *Server) RegisterHandler(httpMethod, relativePath string, handlers ...gin.HandlerFunc) {
	s.Handle(httpMethod, relativePath, handlers...)
}
