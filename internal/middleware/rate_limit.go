package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimit 限流中间件
func RateLimit() gin.HandlerFunc {
	type client struct {
		count  int
		expire time.Time
	}

	clients := make(map[string]*client)
	var mu sync.Mutex

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		defer mu.Unlock()

		if cl, exists := clients[ip]; exists {
			if time.Now().After(cl.expire) {
				// 限流窗口过期，重置计数
				cl.count = 0
				cl.expire = time.Now().Add(time.Minute)
			}

			cl.count++
			if cl.count > 100 { // 每分钟最多100次请求
				c.JSON(429, gin.H{
					"code":  429,
					"error": "请求过于频繁，请稍后重试",
				})
				c.Abort()
				return
			}
		} else {
			// 新客户端
			clients[ip] = &client{
				count:  1,
				expire: time.Now().Add(time.Minute),
			}
		}

		c.Next()
	}
}