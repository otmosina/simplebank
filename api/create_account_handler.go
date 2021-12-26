package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/otmosina/simplebank/db/sqlc"
)

type CreateAccountsRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD RUB IDR"`
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
