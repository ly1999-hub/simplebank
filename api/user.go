package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/ly1999-hub/simplebank/sqlc"
	"github.com/ly1999-hub/simplebank/util"
)

type CreateUserRequest struct {
	Username string `json:"username"  binding:"required"`
	Password string `json:"password"  binding:"required"`
	FullName string `json:"full_name"  binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type UserResponse struct {
	Username         string    `json:"username"`
	FullName         string    `json:"fullname"`
	Email            string    `json:"email"`
	CreatedAt        time.Time `json:"created_at"`
	ChangePasswordAt time.Time `json:"change_password_at"`
}

func (server *Server) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest

	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			fmt.Println(pqErr)
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, ErrorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	res := newUserResponse(user)
	ctx.JSON(http.StatusOK, res)
}

func newUserResponse(user db.User) UserResponse {
	return UserResponse{
		Username:         user.Username,
		FullName:         user.FullName,
		Email:            user.Email,
		ChangePasswordAt: user.ChangePasswordAt.Time,
		CreatedAt:        user.CreatedAt.Time,
	}
}

type userLoginRequest struct {
	Username string `json:"username"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type userLoginResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var userLoginRequest userLoginRequest

	if err := ctx.Bind(&userLoginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
	}

	user, err := server.store.GetUser(ctx, userLoginRequest.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	err = util.CheckPassword(userLoginRequest.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
		return
	}
	accessToken, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}
	rsp := userLoginResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)
}
