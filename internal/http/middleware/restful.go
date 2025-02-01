package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"go_shurtiner/internal/adapter"
)

func RestfulParamsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		pagination := adapter.NewPagination(c.Request)
		ctx := c.Request.Context()

		ctx = context.WithValue(ctx, adapter.Pagination, pagination)

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
