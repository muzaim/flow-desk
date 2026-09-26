package middleware

import (
	"bytes"
	"net/http"
	"time"

	"flow-desk/repository"
	"flow-desk/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func IdempotencyMiddleware(idempotencyRepo repository.IdempotencyRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idempotencyKey := c.GetHeader("Idempotency-Key")

		if idempotencyKey == "" {
			appErr := utils.NewBadRequestError("MISSING_IDEMPOTENCY_KEY", "Header 'Idempotency-Key' wajib diisi untuk membuat task", nil)
			utils.RespondWithError(c, appErr)
			c.Abort()
			return
		}

		if _, err := uuid.Parse(idempotencyKey); err != nil {
			appErr := utils.NewBadRequestError("INVALID_IDEMPOTENCY_KEY", "Format 'Idempotency-Key' tidak valid. Harus berupa UUID v4", err)
			utils.RespondWithError(c, appErr)
			c.Abort()
			return
		}

		userIDVal, exists := c.Get("userID")
		if !exists {
			appErr := utils.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized", nil)
			utils.RespondWithError(c, appErr)
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
