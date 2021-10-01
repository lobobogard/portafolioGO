package model

import "gorm.io/gorm"

type FrontEnd struct {
	gorm.Model
	PerfilID      uint `gorm:"not null;" json:"PerfilID,omitempty"`
	CatfrontendID uint `gorm:"not null;" json:"catfrontendID,omitempty"`
}
