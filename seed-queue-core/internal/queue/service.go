package queue

import (
	"errors"
	"fmt"
	"seed-queue-core/internal/database"
	"seed-queue-core/internal/models"
	"seed-queue-core/internal/pricing"
	"time"
)

type Service struct {
	pricingSvc *pricing.Service
}

func NewService(pricingSvc *pricing.Service) *Service {
	return &Service{pricingSvc: pricingSvc}
}

type CreateOrderRequest struct {
	FarmerName  string           `json:"farmer_name" binding:"required"`
	FarmerPhone string           `json:"farmer_phone"`
	SeedType    models.SeedType  `json:"seed_type" binding:"required"`
	GrossWeight float64          `json:"gross_weight" binding:"gt=0"`
	NetWeight   float64          `json:"net_weight"`
	CakeTaken   bool             `json:"cake_taken"`
	Remark      string           `json:"remark"`
}

func (s *Service) CreateOrder(req *CreateOrderRequest) (*models.QueueOrder, error) {
	calcResult, err := s.pricingSvc.Calculate(req.SeedType, req.GrossWeight, req.CakeTaken)
	if err != nil {
		return nil, errors.New("计费计算失败：" + err.Error())
	}

	farmer := models.Farmer{
		Name:  req.FarmerName,
		Phone: req.FarmerPhone,
	}
	if err := database.DB.Where("name = ? AND phone = ?", farmer.Name, farmer.Phone).FirstOrCreate(&farmer).Error; err != nil {
		return nil, errors.New("农户信息保存失败")
	}

	queueNumber := s.generateQueueNumber()

	order := &models.QueueOrder{
		QueueNumber:   queueNumber,
		FarmerID:      farmer.ID,
		Farmer:        farmer,
		SeedType:      req.SeedType,
		GrossWeight:   req.GrossWeight,
		NetWeight:     req.NetWeight,
		CakeTaken:     req.CakeTaken,
		ProcessingFee: calcResult.ProcessingFee,
		ExpectedOil:   calcResult.ExpectedOil,
		ActualOil:     0,
		PickupStatus:  models.PickupStatusPending,
		PickedUpOil:   0,
		Status:        models.QueueStatusWaiting,
		Remark:        req.Remark,
	}

	if err := database.DB.Create(order).Error; err != nil {
		return nil, errors.New("订单创建失败")
	}
	return order, nil
}

func (s *Service) generateQueueNumber() string {
	now := time.Now()
	var count int64
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	database.DB.Model(&models.QueueOrder{}).Where("created_at >= ?", todayStart).Count(&count)
	return fmt.Sprintf("%s%03d", now.Format("0102"), count+1)
}

func (s *Service) ListOrders(status models.QueueStatus, todayOnly bool) ([]models.QueueOrder, error) {
	var orders []models.QueueOrder
	query := database.DB.Preload("Farmer").Order("created_at asc")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if todayOnly {
		now := time.Now()
		todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		query = query.Where("created_at >= ?", todayStart)
	}
	result := query.Find(&orders)
	return orders, result.Error
}

func (s *Service) GetOrder(id uint) (*models.QueueOrder, error) {
	var order models.QueueOrder
	result := database.DB.Preload("Farmer").First(&order, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &order, nil
}

func (s *Service) CallNext() (*models.QueueOrder, error) {
	var order models.QueueOrder
	result := database.DB.Preload("Farmer").
		Where("status = ?", models.QueueStatusWaiting).
		Order("created_at asc").
		First(&order)
	if result.Error != nil {
		return nil, errors.New("no waiting orders")
	}
	order.Status = models.QueueStatusProcessing
	database.DB.Save(&order)
	return &order, nil
}

func (s *Service) CallSpecific(id uint) (*models.QueueOrder, error) {
	order, err := s.GetOrder(id)
	if err != nil {
		return nil, err
	}
	if order.Status != models.QueueStatusWaiting {
		return nil, errors.New("order is not waiting")
	}
	order.Status = models.QueueStatusProcessing
	database.DB.Save(order)
	return order, nil
}

func (s *Service) CompleteOrder(id uint, actualOil float64) (*models.QueueOrder, error) {
	order, err := s.GetOrder(id)
	if err != nil {
		return nil, err
	}
	if actualOil <= 0 {
		return nil, errors.New("actual oil must be positive")
	}
	now := time.Now()
	order.ActualOil = actualOil
	order.Status = models.QueueStatusCompleted
	order.CompletedAt = &now
	database.DB.Save(order)
	return order, nil
}

func (s *Service) CancelOrder(id uint) error {
	return database.DB.Model(&models.QueueOrder{}).
		Where("id = ?", id).
		Update("status", models.QueueStatusCancelled).Error
}

func (s *Service) PickupOil(id uint, amount float64) (*models.QueueOrder, error) {
	order, err := s.GetOrder(id)
	if err != nil {
		return nil, err
	}
	if order.Status != models.QueueStatusCompleted {
		return nil, errors.New("order not completed yet")
	}
	if order.PickedUpOil+amount > order.ActualOil+0.001 {
		return nil, errors.New("pickup amount exceeds actual oil")
	}
	order.PickedUpOil += amount
	if order.PickedUpOil+0.001 >= order.ActualOil {
		order.PickupStatus = models.PickupStatusDone
	} else {
		order.PickupStatus = models.PickupStatusPartial
	}
	database.DB.Save(order)
	return order, nil
}
