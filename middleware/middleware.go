package middleware

import (
	"net/http"
	"strings"

	"github.com/danielboakye/go-xm/config"
	"github.com/danielboakye/go-xm/helpers"
	"github.com/danielboakye/go-xm/helpers/consts"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderName   = "Authorization"
	authorizationHeaderScheme = "Bearer"
)

func NewRouteFilter(cfg config.Configurations) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		claims, err := getJWTClaimsFromHTTPRequest(ctx, cfg)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": helpers.ErrUnauthorized.Error()})
			ctx.Abort()
			return
		}
		ctx = setJWTClaimsInContext(ctx, claims)

		ctx.Next()
	}
}

func getJWTClaimsFromHTTPRequest(ctx *gin.Context, cfg config.Configurations) (
	claims *helpers.JWTClaims, err error,
) {
	authorizationToken := ctx.GetHeader(authorizationHeaderName)
	clientToken := strings.TrimSpace(strings.Replace(authorizationToken, authorizationHeaderScheme, "", 1))

	if clientToken == "" {
		err = helpers.ErrInvalidToken
		return
	}

	claims, err = helpers.ValidateAccessToken(clientToken, cfg)
	return
}

func setJWTClaimsInContext(ctx *gin.Context, claims *helpers.JWTClaims) *gin.Context {
	ctx.Set(consts.USER_ID, claims.UserID)
	return ctx
}
