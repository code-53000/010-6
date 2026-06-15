package handlers

import (
	"net/http"
	"seed-queue-core/internal/models"
	"seed-queue-core/internal/pricing"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PricingHandler struct {
	svc *pricing.Service
}

func NewPricingHandler(svc *pricing.Service) *PricingHandler {
	return &PricingHandler{svc: svc}
}

func (h *PricingHandler) List(c *gin.Context) {
	configs, err := h.svc.ListPriceConfigs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, configs)
}

func (h *PricingHandler) Update(c *gin.Context) {
	var cfg models.PriceConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.UpdatePriceConfig(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cfg)
}

func (h *PricingHandler) Calculate(c *gin.Context) {
	seedType := models.SeedType(c.Query("seed_type"))
	grossWeight, _ := strconv.ParseFloat(c.Query("gross_weight"), 64)
	cakeTaken, _ := strconv.ParseBool(c.Query("cake_taken"))

	result, err := h.svc.Calculate(seedType, grossWeight, cakeTaken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
