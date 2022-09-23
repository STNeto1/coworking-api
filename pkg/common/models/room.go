package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Room struct {
	ID          string         `sql:"type:uuid;primary_key" json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Schedules   []Schedule     `json:"schedules,omitempty"`
	BuildingID  string         `json:"-"`
	Building    *Building      `json:"building,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (model *Room) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New()
	model.ID = id.String()
	return
}
