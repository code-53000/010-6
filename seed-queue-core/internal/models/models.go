package models

import (
	"time"
)

type SeedType string

const (
	SeedTypePeanut SeedType = "peanut"
	SeedTypeRapeseed SeedType = "rapeseed"
)

type QueueStatus string

const (
	QueueStatusWaiting   QueueStatus = "waiting"
	QueueStatusProcessing QueueStatus = "processing"
	QueueStatusCompleted QueueStatus = "completed"
	QueueStatusCancelled QueueStatus = "cancelled"
)

type OilPickupStatus string

const (
	PickupStatusPending    OilPickupStatus = "pending"
	PickupStatusPartial  OilPickupStatus = "partial"
	PickupStatusDone     OilPickupStatus = "done"
)

type Farmer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Phone     string    `gorm:"size:20" json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

type PriceConfig struct {
	ID             SeedType  `gorm:"primaryKey;size:50" json:"seed_type"`
	PricePerKg     float64 `gorm:"not null" json:"price_per_kg"`
	OilRate        float64 `gorm:"not null" json:"oil_rate"`
	ProcessingFee float64 `gorm:"not null" json:"processing_fee"`
	CakeTakeFee    float64 `gorm:"not null;default:0" json:"cake_take_fee"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type QueueOrder struct {
	ID              uint            `gorm:"primaryKey" json:"id"`
	QueueNumber     string          `gorm:"size:20;uniqueIndex;not null" json:"queue_number"`
	FarmerID        uint            `json:"farmer_id"`
	Farmer          Farmer          `gorm:"foreignKey:FarmerID" json:"farmer"`
	SeedType        SeedType        `gorm:"size:50;not null" json:"seed_type"`
	GrossWeight     float64         `gorm:"not null" json:"gross_weight"`
	NetWeight       float64         `json:"net_weight"`
	CakeTaken       bool            `gorm:"default:false" json:"cake_taken"`
	ProcessingFee   float64         `json:"processing_fee"`
	ExpectedOil    float64         `json:"expected_oil"`
	ActualOil      float64         `json:"actual_oil"`
	PickupStatus   OilPickupStatus `gorm:"size:20;default:pending" json:"pickup_status"`
	PickedUpOil    float64         `gorm:"default:0" json:"picked_up_oil"`
	Status          QueueStatus     `gorm:"size:20;default:waiting" json:"status"`
	Remark         string          `gorm:"size:500" json:"remark"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	CompletedAt    *time.Time     `json:"completed_at"`
}

type OilStorage struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	OrderID      uint      `gorm:"index;not null" json:"order_id"`
	Order        QueueOrder `gorm:"foreignKey:OrderID" json:"order"`
	FarmerID     uint      `json:"farmer_id"`
	Farmer       Farmer    `gorm:"foreignKey:FarmerID" json:"farmer"`
	BarrelCount   int       `gorm:"not null" json:"barrel_count"`
	OilPerBarrel float64   `gorm:"not null" json:"oil_per_barrel"`
	TotalOil     float64   `gorm:"not null" json:"total_oil"`
	StoredDays   int       `gorm:"default:0" json:"stored_days"`
	StorageFee   float64   `gorm:"default:0" json:"storage_fee"`
	IsPickedUp   bool      `gorm:"default:false" json:"is_picked_up"`
	PickedUpAt   *time.Time `json:"picked_up_at"`
	CreatedAt    time.Time `json:"created_at"`
	Remark       string    `gorm:"size:500" json:"remark"`
}
