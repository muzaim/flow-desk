package middleware

import (
	"bytes"
	"flow-desk/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func IdempotencyMiddleware(idempotencyRepo repository.IdempotencyRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idempotencyKey := c.GetHeader("Idempotency-Key")

		if idempotencyKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Header 'Idempotency-Key' wajib diisi untuk membuat task",
			})
			c.Abort()
			return
		}

		if _, err := uuid.Parse(idempotencyKey); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Format 'Idempotency-Key' tidak valid. Harus berupa UUID v4",
			})
			c.Abort()
			return
		}

		userIDVal, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		userID := userIDVal.(uint)

		cachedData, err := idempotencyRepo.Get(c.Request.Context(), idempotencyKey, userID)
		if err == nil && cachedData != nil {
			c.Header("Content-Type", "application/json")
			c.Header("X-Cache-Lookup", "HIT (Redis)")
			c.String(cachedData.ResponseCode, cachedData.ResponseBody)
			c.Abort()
			return
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		statusCode := c.Writer.Status()
		if statusCode >= 200 && statusCode < 300 {
			dataToSave := &repository.IdempotencyData{
				ResponseCode: statusCode,
				ResponseBody: blw.body.String(),
			}
			_ = idempotencyRepo.Set(c.Request.Context(), idempotencyKey, userID, dataToSave, 24*time.Hour)
		}
	}
}
