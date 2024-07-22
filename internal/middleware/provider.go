package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"net/http"
)

func GothProvider(c *gin.Context) {
	gothic.GetProviderName = func(req *http.Request) (string, error) {
		provider := c.Param("provider")
		if provider == "" {
			return "", errors.New("no provider specified")
		}
		return provider, nil
	}
	c.Next()
}
