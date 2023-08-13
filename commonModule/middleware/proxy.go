package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http/httputil"
	"net/url"
)

// Proxy 代理中间件
func Proxy(target *url.URL) gin.HandlerFunc {
	return func(c *gin.Context) {
		proxy := httputil.NewSingleHostReverseProxy(target)
		c.Request.Host = target.Host
		c.Request.URL.Path = c.Param("path")
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
