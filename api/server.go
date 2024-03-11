package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/ly1999-hub/simplebank/sqlc"
	"github.com/ly1999-hub/simplebank/token"
	"github.com/ly1999-hub/simplebank/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can not create token maker")
	}
	server := Server{config: config, store: store, tokenMaker: tokenMaker}
	server.setUpRouter()
	return &server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func ErrorResponse(err error) gin.H {
	return gin.H{"error ": err.Error()}
}

func (server *Server) setUpRouter() {
	router := gin.Default()
	router.POST("/users", server.CreateUser)
	router.POST("/users/login-by-username", server.loginUser)

	authRouter := router.Group("/auth").Use(authMiddleware(server.tokenMaker))
	authRouter.POST("/accounts", server.createAccount)
	authRouter.GET("/accounts/:id", server.getAccount)
	authRouter.GET("/accounts", server.getListAccount)
	authRouter.POST("/transfers", server.CreateTransfer)

	server.router = router
}
