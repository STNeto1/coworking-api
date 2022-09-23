package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	RoomPricingTypeShift  = "shift"
	RoomPricingTypeDaily  = "daily"
	RoomPricingTypeWeekly = "weekly"
)

type RoomPricing struct {
	ID        string         `sql:"type:uuid;primary_key" json:"id"`
	Price     float32        `json:"price"`
	Type      string         `json:"type"`
	RoomID    string         `json:"-"`
	Room      *Room          `json:"room,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (model *RoomPricing) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New()
	model.ID = id.String()
	return
}
