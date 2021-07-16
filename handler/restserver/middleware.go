package restserver

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

func ipFilter(ipNets []*net.IPNet) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := net.ParseIP(c.ClientIP())
		for _, ip := range ipNets {
			if ip.Contains(clientIP) {
				c.Next()
				return
			}
		}
		c.AbortWithStatus(http.StatusForbidden)
	}
}

func setApiKey(c *gin.Context) {
	apiKey := c.Request.Header.Get("Authorization")
	ctx := c.Request.Context()
	md := metadata.New(map[string]string{"authorization": apiKey})
	ctx = metadata.NewOutgoingContext(ctx, md)
	*c.Request = *c.Request.WithContext(ctx)
	c.Next()
}
