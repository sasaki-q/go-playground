package src

import (
	db "dbapp/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := Server{store: store}
	router := gin.Default()

	server.router = router

	router.GET("/account/:id", server.getAccount)
	router.POST("/account", server.createAccount)

	return &server
}

func (server *Server) StartServer(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
