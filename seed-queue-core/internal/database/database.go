package database

import (
	"seed-queue-core/internal/config"
	"seed-queue-core/internal/models"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(cfg *config.Config) error {
	var err error
	for i := 0; i < 30; i++ {
		DB, err = gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return err
	}

	if err = DB.AutoMigrate(
		&models.Farmer{},
		&models.PriceConfig{},
		&models.QueueOrder{},
		&models.OilStorage{},
	); err != nil {
		return err
	}

	seedPriceConfigs()
	return nil
}

func seedPriceConfigs() {
	defaults := []models.PriceConfig{
		{ID: models.SeedTypePeanut, PricePerKg: 30.00, OilRate: 45.0, ProcessingFee: 50.00, CakeTakeFee: 10.00, UpdatedAt: time.Now()},
		{ID: models.SeedTypeRapeseed, PricePerKg: 25.00, OilRate: 38.0, ProcessingFee: 50.00, CakeTakeFee: 8.00, UpdatedAt: time.Now()},
	}
	for _, d := range defaults {
		var count int64
		DB.Model(&models.PriceConfig{}).Where("id = ?", d.ID).Count(&count)
		if count == 0 {
			DB.Create(&d)
		}
	}
}
