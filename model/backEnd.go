package model

import "gorm.io/gorm"

type BackEnd struct {
	gorm.Model
	PerfilID     uint   `gorm:"not null;" json:"PerfilID,omitempty"`
	Perfil       Perfil `gorm:"foreignKey:PerfilID;"`
	CatbackendID uint   `gorm:"not null;" json:"catbackendID,omitempty"`
	Status       bool   `gorm:"not null;" json:"statusServer,omitempty"`
}
