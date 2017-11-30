package web

import (
	"github.com/gin-gonic/gin"
	"github.com/remyr/goauth-microservice/utils"
	"github.com/remyr/goauth-microservice/controllers"
)

type Server struct {
	*gin.Engine
}

func NewServer(dba utils.DatabaseAccessor) *Server {
	router := gin.Default()
	server := &Server{router}

	userController := controllers.NewUserController(dba)
	userController.Register(router)

	return server
}
