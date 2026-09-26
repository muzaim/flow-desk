package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"flow-desk/utils"

	"github.com/gin-gonic/gin"
)

func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC RECOVERED] Error: %v", r)

				c.JSON(http.StatusInternalServerError, utils.ErrorResponse{
					Status:    false,
					Code:      "PANIC_SERVER_ERROR",
					Message:   fmt.Sprintf("Internal server panic: %v", r),
					Timestamp: time.Now(),
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
