package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/ly1999-hub/simplebank/sqlc"
	"github.com/ly1999-hub/simplebank/token"
)

type CreateAccountRequest struct {
	Owner   string `json:"owner" binding:"required"`
	Curency string `json:"curency" binding:"required,oneof=USD UER"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner:   authPayload.Username,
		Balance: 0,
		Curency: req.Curency,
	}

	account, err := server.store.CreateAccount(ctx, arg)

	if err != nil {
		fmt.Println(err.(*pq.Error).Code.Name())
		if pqErr, ok := err.(*pq.Error); ok {
			fmt.Println(pqErr)
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, ErrorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, account)

}

type getAccountRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest

	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
	account, err := server.store.GetAccount(ctx, int64(req.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err := errors.New("account dont belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
	}
	ctx.JSON(http.StatusOK, account)
}

type getListAccountRequest struct {
	Limit int32 `form:"limit" binding:"required,min=10,max=1000"`
	Page  int32 `form:"page" binding:"required,min=1"`
}

func (server *Server) getListAccount(ctx *gin.Context) {
	var req getListAccountRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.Limit,
		Offset: (req.Page - 1) * req.Limit,
	}
	accounts, err := server.store.ListAccounts(context.Background(), arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}
