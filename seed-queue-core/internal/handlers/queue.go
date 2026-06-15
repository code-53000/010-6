package handlers

import (
	"net/http"
	"seed-queue-core/internal/models"
	"seed-queue-core/internal/queue"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QueueHandler struct {
	svc *queue.Service
}

func NewQueueHandler(svc *queue.Service) *QueueHandler {
	return &QueueHandler{svc: svc}
}

func (h *QueueHandler) Create(c *gin.Context) {
	var req queue.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order, err := h.svc.CreateOrder(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, order)
}

func (h *QueueHandler) List(c *gin.Context) {
	status := models.QueueStatus(c.Query("status"))
	todayOnly, _ := strconv.ParseBool(c.Query("today_only"))
	orders, err := h.svc.ListOrders(status, todayOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *QueueHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	order, err := h.svc.GetOrder(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *QueueHandler) CallNext(c *gin.Context) {
	order, err := h.svc.CallNext()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *QueueHandler) CallSpecific(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	order, err := h.svc.CallSpecific(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

type CompleteRequest struct {
	ActualOil float64 `json:"actual_oil" binding:"required,gt=0"`
}

func (h *QueueHandler) Complete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req CompleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order, err := h.svc.CompleteOrder(uint(id), req.ActualOil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *QueueHandler) Cancel(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.CancelOrder(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cancelled"})
}

type PickupRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

func (h *QueueHandler) PickupOil(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req PickupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order, err := h.svc.PickupOil(uint(id), req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}
