package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Address struct {
	ID            string `sql:"type:uuid;primary_key" json:"id"`
	PostalCode    string `json:"postal_code"`
	State         string `json:"state"`
	City          string `json:"city"`
	Street        string `json:"street"`
	Number        string `json:"number"`
	GoogleMapsURL string `json:"google_maps_url"`
}

func (model *Address) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New()
	model.ID = id.String()
	return
}
