package pricing

import (
	"errors"
	"seed-queue-core/internal/database"
	"seed-queue-core/internal/models"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

type CalculateResult struct {
	ProcessingFee float64 `json:"processing_fee"`
	ExpectedOil   float64 `json:"expected_oil"`
}

func (s *Service) GetPriceConfig(seedType models.SeedType) (*models.PriceConfig, error) {
	var cfg models.PriceConfig
	result := database.DB.Where("id = ?", seedType).First(&cfg)
	if result.Error != nil {
		return nil, result.Error
	}
	return &cfg, nil
}

func (s *Service) ListPriceConfigs() ([]models.PriceConfig, error) {
	var configs []models.PriceConfig
	result := database.DB.Find(&configs)
	return configs, result.Error
}

func (s *Service) UpdatePriceConfig(cfg *models.PriceConfig) error {
	if cfg.ID == "" {
		return errors.New("seed type is required")
	}
	if cfg.PricePerKg < 0 || cfg.ProcessingFee < 0 || cfg.OilRate <= 0 {
		return errors.New("invalid price values")
	}
	return database.DB.Save(cfg).Error
}

func (s *Service) Calculate(seedType models.SeedType, grossWeight float64, cakeTaken bool) (*CalculateResult, error) {
	cfg, err := s.GetPriceConfig(seedType)
	if err != nil {
		return nil, err
	}

	expectedOil := grossWeight * cfg.OilRate / 100.0

	processingFee := grossWeight*cfg.PricePerKg/100.0 + cfg.ProcessingFee
	if cakeTaken {
		processingFee += grossWeight * cfg.CakeTakeFee / 100.0
	}

	return &CalculateResult{
		ProcessingFee: roundTo2(processingFee),
		ExpectedOil:   roundTo2(expectedOil),
	}, nil
}

func roundTo2(v float64) float64 {
	return float64(int(v*100+0.5)) / 100.0
}
