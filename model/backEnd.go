package model

import "gorm.io/gorm"

type BackEnd struct {
	gorm.Model
	PerfilID     uint `gorm:"not null;" json:"PerfilID,omitempty"`
	CatbackendID uint `gorm:"not null;" json:"catbackendID,omitempty"`
}
