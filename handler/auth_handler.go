package handler

import (
	"net/http"

	"flow-desk/dto"
	"flow-desk/service"
	"flow-desk/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := utils.NewBadRequestError("INVALID_INPUT", "Format payload request tidak valid", err)
		utils.RespondWithError(c, appErr)
		return
	}

	res, err := h.authService.Register(req)
	if err != nil {
		utils.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registrasi berhasil",
		"data":    res,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := utils.NewBadRequestError("INVALID_INPUT", "Format payload request tidak valid", err)
		utils.RespondWithError(c, appErr)
		return
	}

	res, err := h.authService.Login(req)
	if err != nil {
		utils.RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"data":    res,
	})
}
