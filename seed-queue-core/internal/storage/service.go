package storage

import (
	"errors"
	"seed-queue-core/internal/database"
	"seed-queue-core/internal/models"
	"time"
)

type Service struct {
	StorageFeePerDay float64
}

func NewService() *Service {
	return &Service{StorageFeePerDay: 0.5}
}

type CreateStorageRequest struct {
	OrderID      uint    `json:"order_id" binding:"required"`
	BarrelCount   int     `json:"barrel_count" binding:"required,gt=0"`
	OilPerBarrel float64 `json:"oil_per_barrel" binding:"required,gt=0"`
	StoredDays   int     `json:"stored_days"`
	Remark       string  `json:"remark"`
}

func (s *Service) CreateStorage(req *CreateStorageRequest) (*models.OilStorage, error) {
	var order models.QueueOrder
	if err := database.DB.First(&order, req.OrderID).Error; err != nil {
		return nil, errors.New("order not found")
	}
	if order.Status != models.QueueStatusCompleted {
		return nil, errors.New("order not completed")
	}

	totalOil := float64(req.BarrelCount) * req.OilPerBarrel
	remainingOil := order.ActualOil - order.PickedUpOil
	if totalOil > remainingOil+0.001 {
		return nil, errors.New("寄存油量超出剩余可取油量")
	}

	storageFee := float64(req.StoredDays) * s.StorageFeePerDay * float64(req.BarrelCount)

	storage := &models.OilStorage{
		OrderID:      req.OrderID,
		FarmerID:     order.FarmerID,
		BarrelCount:   req.BarrelCount,
		OilPerBarrel: req.OilPerBarrel,
		TotalOil:     totalOil,
		StoredDays:   req.StoredDays,
		StorageFee:   storageFee,
		IsPickedUp:   false,
		CreatedAt:    time.Now(),
		Remark:       req.Remark,
	}

	tx := database.DB.Begin()
	if err := tx.Create(storage).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	order.PickedUpOil += totalOil
	if order.PickedUpOil+0.001 >= order.ActualOil {
		order.PickupStatus = models.PickupStatusDone
	} else {
		order.PickupStatus = models.PickupStatusPartial
	}
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return storage, nil
}

func (s *Service) ListStorage(isPickedUp *bool) ([]models.OilStorage, error) {
	var storages []models.OilStorage
	query := database.DB.Preload("Order").Preload("Farmer").Order("created_at desc")
	if isPickedUp != nil {
		query = query.Where("is_picked_up = ?", *isPickedUp)
	}
	result := query.Find(&storages)
	return storages, result.Error
}

func (s *Service) PickupStorage(id uint) (*models.OilStorage, error) {
	var storage models.OilStorage
	if err := database.DB.First(&storage, id).Error; err != nil {
		return nil, err
	}
	if storage.IsPickedUp {
		return nil, errors.New("already picked up")
	}
	now := time.Now()
	storage.IsPickedUp = true
	storage.PickedUpAt = &now

	extraDays := int(now.Sub(storage.CreatedAt).Hours()/24) - storage.StoredDays
	if extraDays > 0 {
		storage.StorageFee += float64(extraDays) * s.StorageFeePerDay * float64(storage.BarrelCount)
	}

	database.DB.Save(&storage)
	return &storage, nil
}

func (s *Service) GetStorage(id uint) (*models.OilStorage, error) {
	var storage models.OilStorage
	result := database.DB.Preload("Order").Preload("Farmer").First(&storage, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &storage, nil
}
