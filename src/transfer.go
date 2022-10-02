package src

import (
	"database/sql"
	db "dbapp/db/sqlc"
	"dbapp/factory"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var (
		req transferRequest
		err error
	)

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	param := db.CreateTransferParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	account, valid := server.validateAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*factory.Payload)
	if account.Owner != payload.Username {
		err := errors.New("account is not authenticated user's account")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, valid = server.validateAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res, err := server.store.TransferTransaction(ctx, param)

	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			errorResponse(err),
		)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (server *Server) validateAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.SelectAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account %d currency mismatch %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	return account, true
}
