package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SecureCookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetSameSite(http.SameSiteNoneMode)
		c.Next()
	}
}
