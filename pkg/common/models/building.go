package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Building struct {
	ID          string         `sql:"type:uuid;primary_key" json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	UserID      string         `json:"-"`
	User        User           `json:"-"`
	AddressID   string         `json:"-"`
	Address     *Address       `json:"address,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (model *Building) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New()
	model.ID = id.String()
	return
}
