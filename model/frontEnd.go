package model

import "gorm.io/gorm"

type FrontEnd struct {
	gorm.Model
	PerfilID      uint   `gorm:"not null;" json:"PerfilID,omitempty"`
	Perfil        Perfil `gorm:"foreignKey:PerfilID;"`
	CatfrontendID uint   `gorm:"not null;" json:"catfrontendID,omitempty"`
	Status        bool   `gorm:"not null;" json:"statusServer,omitempty"`
}
