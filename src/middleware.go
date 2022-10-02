package src

import (
	"dbapp/factory"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var authorizationPayloadKey = "authorization_payload"

func authMiddleware(factory factory.Factory) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("authorization")
		if len(token) == 0 {
			err := errors.New("authorization token is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(token)
		if len(fields) < 2 {
			err := errors.New("invalid token format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		if fields[0] != "Bearer" {
			err := fmt.Errorf("unsupported authorization type %s", fields[0])
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		payload, err := factory.VerifyToken(fields[1])

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()

	}
}
