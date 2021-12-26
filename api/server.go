package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/otmosina/simplebank/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

type IndexAccountsRequest struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5"`
}

func NewServer(db db.Store) *Server {
	server := &Server{store: db}
	router := gin.Default()

	router.GET("/accounts", server.indexAccounts)
	router.GET("/account/:id", server.getAccount)
	router.POST("/account/:id", server.deleteAccount)
	router.POST("/accounts", server.createAccounts)
	// TODO add some routes

	server.router = router
	return server

}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) indexAccounts(ctx *gin.Context) {
	var req IndexAccountsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	accounts, err := server.store.ListAccounts(ctx, db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}
