package middleware

import (
	"log"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"github.com/kodokbakar/go-ecommerce-api/internal/response"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Printf("panic recovered: %v\n%s", recovered, string(debug.Stack()))

				if !c.Writer.Written() {
					response.InternalServerError(c, "internal server error", nil)
				}

				c.Abort()
			}
		}()

		c.Next()
	}
}
