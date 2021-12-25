package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/otmosina/simplebank/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

type CreateAccountsRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD RUB IDR"`
}

type IndexAccountsRequest struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5"`
}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type DeleteAccountRequest GetAccountRequest

func NewServer(db db.Store) *Server {
	server := &Server{store: db}
	router := gin.Default()

	// router.POST("/accounts", createAccounts())
	// router.GET("/accounts", func(ctx *gin.Context) {
	// 	response := gin.H{
	// 		"status": "ok",
	// 	}
	// 	ctx.JSON(http.StatusOK, response)
	// })

	router.GET("/accounts", server.indexAccounts)
	router.GET("/account/:id", server.getAccount)
	router.POST("/accounts/:id", server.deleteAccount)
	router.POST("/accounts", server.createAccounts)
	// TODO add some routes

	server.router = router
	return server

}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req DeleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, ErrorResponse(err))
		return
	}
	err = server.store.DeleteAccount(ctx, account.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse())
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req GetAccountRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	// account = db.Account{}
	ctx.JSON(http.StatusOK, account)
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

func (server *Server) createAccounts(ctx *gin.Context) {
	// var err error
	// var account db.Account
	var req CreateAccountsRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	account, err := server.store.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, account)
}
