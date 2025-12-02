package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Auth JWT认证中间件
func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// 1. 尝试从 Header 获取
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}

		// 2. 尝试从 Query 参数获取 (用于 EventSource 等不支持自定义 Header 的场景)
		if tokenString == "" {
			tokenString = c.Query("token")
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":  401,
				"error": "缺少认证信息",
			})
			c.Abort()
			return
		}

		// 解析JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法，防止算法替换攻击
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":  401,
				"error": "无效的token",
			})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":  401,
				"error": "token已失效",
			})
			c.Abort()
			return
		}

		// 获取用户ID
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":  401,
				"error": "无效的token声明",
			})
			c.Abort()
			return
		}

		// 将用户ID设置到上下文中
		if userID, exists := claims["userId"]; exists {
			c.Set("userId", uint(userID.(float64)))
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":  401,
				"error": "无效的token: 缺少用户ID",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserID 从上下文中获取用户ID
// 此时中间件已确保用户已认证且ID存在，因此直接返回ID即可
func GetUserID(c *gin.Context) uint {
	return c.GetUint("userId")
}
