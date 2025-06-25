package middleware

import (
	"go_shurtiner/internal/app/authentication"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:generate mockery --name=AuthenticationMiddleware --output=../../../test/mocks --outpkg=mocks --filename=mock_authentication_middleware.go
func AuthenticationMiddleware(auth authentication.Authentication) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := auth.Authenticate(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, auth.UnauthorizedResponse(err))
			return
		}
		c.Set(authentication.User, user)
		c.Next()
	}
}
