package handler

import (
	"net/http"

	"flow-desk/service"

	"github.com/gin-gonic/gin"
)

type RegionHandler struct {
	regionService service.RegionService
}

func NewRegionHandler(regionService service.RegionService) *RegionHandler {
	return &RegionHandler{regionService: regionService}
}

func (h *RegionHandler) GetProvinces(c *gin.Context) {
	provinces, err := h.regionService.GetProvinces()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data provinsi",
		"data":    provinces,
	})
}
