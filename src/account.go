package src

import (
	db "dbapp/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var (
		req createAccountRequest
		err error
	)

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			errorResponse(err),
		)
		return
	}

	param := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, param)

	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			errorResponse(err),
		)
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var (
		req getAccountRequest
		err error
	)

	err = ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			errorResponse(err),
		)
		return
	}

	account, err := server.store.SelectAccount(ctx, req.ID)

	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			errorResponse(err),
		)
		return
	}

	ctx.JSON(http.StatusOK, account)
}