package middleware

import (
	"github.com/gin-gonic/gin"
)

// MaxAllowed limit n connection
func MaxAllowed(n int) gin.HandlerFunc {
	sem := make(chan struct{}, n)
	acquire := func() { sem <- struct{}{} }
	release := func() { <-sem }
	return func(c *gin.Context) {
		acquire() // before request
		c.Next()
		release() // after request
	}
}
