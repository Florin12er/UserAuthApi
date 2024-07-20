package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	lastSeen time.Time
	count    int
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
)

func RateLimiter(limit int, windowSize time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		mu.Lock()
		v, exists := visitors[ip]
		if !exists {
			visitors[ip] = &visitor{lastSeen: time.Now()}
			v = visitors[ip]
		}
		
		if time.Since(v.lastSeen) > windowSize {
			v.count = 0
			v.lastSeen = time.Now()
		}

		if v.count >= limit {
			mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		v.count++
		mu.Unlock()

		c.Next()
	}
}

