package middleware

import (
	"net/http"
	"strings"

	"github.com/Signal-zxh/signal-zxh/model"
	"github.com/Signal-zxh/signal-zxh/utils"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		auth := c.GetHeader("Authorization")

		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.Fail("no token"))
			return
		}
		// 拆Bearer
		parts := strings.Split(auth, " ")

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.Fail("bad format"))
			return
		}
		// 验证token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.Fail("invalid token"))
			return
		}
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
