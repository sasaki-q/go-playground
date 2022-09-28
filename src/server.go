package src

import (
	db "dbapp/db/sqlc"
	"dbapp/factory"
	"dbapp/utils"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	factory factory.Factory
	config  utils.Config
	store   *db.Store
	router  *gin.Engine
}

func NewServer(config utils.Config, store *db.Store) (*Server, error) {
	factory, err := factory.GeneratePasetoFactory(config.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("ERROR: cannot create token factory: %w", err)
	}
	router := gin.Default()
	server := Server{
		store:   store,
		router:  router,
		factory: factory,
		config:  config,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validateCurrency)
	}

	server.setupRouter()
	return &server, nil
}

func (server *Server) setupRouter() {
	router := *gin.Default()

	router.GET("/account/:id", server.getAccount)
	router.POST("/account", server.createAccount)

	router.POST("/transfer", server.createTransfer)
	router.POST("/user", server.createUser)
	router.POST("/user/login", server.login)
}

func (server *Server) StartServer(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
