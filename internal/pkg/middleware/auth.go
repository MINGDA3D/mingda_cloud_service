package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"mingda_cloud_service/internal/pkg/errors"
)

// AuthRequired 认证中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(401, errors.New(errors.ErrUnauthorized, "missing authorization header"))
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(401, errors.New(errors.ErrUnauthorized, "invalid authorization format"))
			return
		}

		// TODO: 验证token
		token := parts[1]
		if !validateToken(token) {
			c.AbortWithStatusJSON(401, errors.New(errors.ErrUnauthorized, "invalid token"))
			return
		}

		c.Next()
	}
}

func validateToken(token string) bool {
	// TODO: 实现token验证
	return true
}
