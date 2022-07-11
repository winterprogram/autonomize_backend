package auth

import (
	"errors"
	"net/http"
	"test/test_app/app/api/middleware/jwt"
	"test/test_app/app/constants"
	"test/test_app/app/controller"
	"test/test_app/app/service/correlation"

	// "test/test_app/app/service/logger"
	"strings"

	"github.com/gin-gonic/gin"
)

//prefix can be changed w.r.t application requirements
var AuthPrefix = "Bearer "
var TokenKey = "Authorization"

// authentication is a middleware that verify JWT token headers
func Authentication(jwt jwt.IJwtService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c := correlation.WithReqContext(ctx)
		token, err := getHeaderToken(ctx)
		if err != nil {
			controller.RespondWithError(ctx, http.StatusUnauthorized, err.Error())
			return
		}

		claims, valid := jwt.VerifyToken(c, token)
		if !valid {
			controller.RespondWithError(ctx, http.StatusUnauthorized, "unauthorized")
			return
		}
		ctx.Set(constants.CtxClaims, claims)
		ctx.Next()
	}
}

func getHeaderToken(ctx *gin.Context) (string, error) {
	header := string(ctx.GetHeader(TokenKey))
	return extractToken(header)
}

func extractToken(header string) (string, error) {
	if strings.HasPrefix(header, AuthPrefix) {
		return header[len(AuthPrefix):], nil
	}
	return "", errors.New("token not found")
}
