package main

import (
	"log"
	"seed-queue-core/internal/config"
	"seed-queue-core/internal/database"
	"seed-queue-core/internal/handlers"
	"seed-queue-core/internal/pricing"
	"seed-queue-core/internal/queue"
	"seed-queue-core/internal/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	if err := database.Init(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connected and migrated")

	pricingSvc := pricing.NewService()
	queueSvc := queue.NewService(pricingSvc)
	storageSvc := storage.NewService()

	pricingH := handlers.NewPricingHandler(pricingSvc)
	queueH := handlers.NewQueueHandler(queueSvc)
	storageH := handlers.NewStorageHandler(storageSvc)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type"},
	}))

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		pricing := api.Group("/pricing")
		{
			pricing.GET("", pricingH.List)
			pricing.PUT("", pricingH.Update)
			pricing.GET("/calculate", pricingH.Calculate)
		}

		queue := api.Group("/queue")
		{
			queue.POST("", queueH.Create)
			queue.GET("", queueH.List)
			queue.GET("/:id", queueH.Get)
			queue.POST("/call-next", queueH.CallNext)
			queue.POST("/:id/call", queueH.CallSpecific)
			queue.POST("/:id/complete", queueH.Complete)
			queue.POST("/:id/cancel", queueH.Cancel)
			queue.POST("/:id/pickup", queueH.PickupOil)
		}

		storage := api.Group("/storage")
		{
			storage.POST("", storageH.Create)
			storage.GET("", storageH.List)
			storage.GET("/:id", storageH.Get)
			storage.POST("/:id/pickup", storageH.Pickup)
		}
	}

	log.Println("Server starting on :" + cfg.ServerPort)
	r.Run(":" + cfg.ServerPort)
}
