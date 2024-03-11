package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ly1999-hub/simplebank/token"
)

const (
	authorizationHearderKey = "authorization"
	authorizationHeaderType = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHearderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationHeaderType {
			err := fmt.Errorf("unsuport authorization header type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
