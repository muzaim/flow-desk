package utils

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RespondWithError(c *gin.Context, err error) {
	var appErr *AppError

	if errors.As(err, &appErr) {
		if appErr.StatusCode >= 500 {
			log.Printf("[SERVER ERROR] Code: %s, Cause: %v", appErr.Code, appErr.Err)
		}

		c.JSON(appErr.StatusCode, ErrorResponse{
			Status:    false,
			Code:      appErr.Code,
			Message:   appErr.Message,
			Timestamp: time.Now(),
		})
		return
	}

	log.Printf("[UNHANDLED ERROR] Cause: %v", err)
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Status:    false,
		Code:      "INTERNAL_SERVER_ERROR",
		Message:   "Terjadi kesalahan internal pada server",
		Timestamp: time.Now(),
	})
}
