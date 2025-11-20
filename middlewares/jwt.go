package middlewares

import (
	"fmt"
	"net/http"
	"webserver/models"
	"webserver/utils"

	"github.com/gin-gonic/gin"
)

// 处理 JWT 认证的中间件
// token 是否合法
func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := utils.TokenValid(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": err.Error()})
			return
		}

		authUser, _ := models.GetUserByID(claims.UserID)
		if authUser == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "用户不存在"})
			return
		}
		c.Set("authUser", authUser)

		c.Next()
	}
}

// CORSMiddleware 中间件处理跨域问题
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		// 获取到请求头中的 Origin
		origin := c.Request.Header.Get("Origin")
		fmt.Println("origin:", origin)
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
