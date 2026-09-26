package middleware

import (
	"net/http"
	"strings"

	"flow-desk/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			appErr := utils.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Header Authorization diperlukan", nil)
			utils.RespondWithError(c, appErr)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			appErr := utils.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Format token harus 'Bearer <token>'", nil)
			utils.RespondWithError(c, appErr)
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			appErr := utils.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Token tidak valid atau sudah kadaluwarsa", err)
			utils.RespondWithError(c, appErr)
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)

		c.Next()
	}
}
