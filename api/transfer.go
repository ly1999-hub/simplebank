package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/ly1999-hub/simplebank/sqlc"
	"github.com/ly1999-hub/simplebank/token"
)

type CreateTransferRequest struct {
	FromAccountID int64  `json:"from_account_id"  binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,min=1"`
	Currency      string `json:"currency" binding:"required,oneof=USD CAD VND EUR"`
}

func (server *Server) CreateTransfer(ctx *gin.Context) {
	var req CreateTransferRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
	fromAccount, invalid := server.validAccount(ctx, req.FromAccountID, req.Currency)
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account dont belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
	}
	if !invalid {
		return
	}
	_, invalid = server.validAccount(ctx, req.FromAccountID, req.Currency)
	if !invalid {
		return
	}
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	transfer, err := server.store.TransferTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, transfer)

}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ErrorResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return account, false
	}
	if account.Curency != currency {
		err := fmt.Errorf("account [%d] currency not mismatch: %s vs %s", account.ID, account.Curency, currency)
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return account, false
	}
	return account, true
}
