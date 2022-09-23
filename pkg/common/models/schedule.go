package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	ScheduleStatusOpen      = "open"
	ScheduleStatusReserved  = "reserved"
	ScheduleStatusConfirmed = "confirmed"
	ScheduleStatusRejected  = "rejected"
)

type Schedule struct {
	ID        string         `sql:"type:uuid;primary_key" json:"id"`
	Price     float32        `json:"Price"`
	From      time.Time      `json:"from"`
	To        time.Time      `json:"to"`
	Status    string         `json:"status"`
	RoomID    string         `json:"-"`
	Room      *Room          `json:"room,omitempty"`
	UserID    sql.NullString `json:"-"`
	User      *User          `json:"user,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (model *Schedule) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New()
	model.ID = id.String()
	return
}
