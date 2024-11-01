package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"waf-detection/internal/services"
)

type FingerprintController struct {
	service *services.FingerprintService
}

func NewFingerprintController(service *services.FingerprintService) *FingerprintController {
	return &FingerprintController{
		service: service,
	}
}

func (con *FingerprintController) CheckIP(c *gin.Context) {
	ctx := c.Request.Context()
	ip := c.Param("ip")

	if ip == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ip must be specified",
		})
		return
	}

	isSuspicious, err := con.service.CheckIP(ctx, ip)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to check ip",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"isSuspicious": isSuspicious,
	})
}
