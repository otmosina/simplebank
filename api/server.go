package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/otmosina/simplebank/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(db db.Store) *Server {
	server := &Server{store: db}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.GET("/accounts", server.indexAccounts)
	router.GET("/account/:id", server.getAccount)
	router.POST("/account/:id", server.deleteAccount)
	router.POST("/accounts", server.createAccounts)
	router.POST("/users", server.createUsers)

	router.POST("/transfers", server.TransferRequest)
	// TODO add some routes

	server.router = router
	return server

}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
