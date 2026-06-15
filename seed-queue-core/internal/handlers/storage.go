package handlers

import (
	"net/http"
	"seed-queue-core/internal/storage"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StorageHandler struct {
	svc *storage.Service
}

func NewStorageHandler(svc *storage.Service) *StorageHandler {
	return &StorageHandler{svc: svc}
}

func (h *StorageHandler) Create(c *gin.Context) {
	var req storage.CreateStorageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s, err := h.svc.CreateStorage(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, s)
}

func (h *StorageHandler) List(c *gin.Context) {
	pickedUpStr := c.Query("is_picked_up")
	var isPickedUp *bool
	if pickedUpStr != "" {
		v := pickedUpStr == "true"
		isPickedUp = &v
	}
	list, err := h.svc.ListStorage(isPickedUp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *StorageHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	s, err := h.svc.GetStorage(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *StorageHandler) Pickup(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	s, err := h.svc.PickupStorage(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}
