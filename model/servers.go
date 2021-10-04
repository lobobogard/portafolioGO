package model

import "gorm.io/gorm"

type Servers struct {
	gorm.Model
	PerfilID    uint   `gorm:"not null;" json:"PerfilID,omitempty"`
	Perfil      Perfil `gorm:"foreignKey:PerfilID;"`
	CatserverID uint   `gorm:"not null;" json:"catserverID,omitempty"`
	Status      bool   `gorm:"not null;" json:"statusServer,omitempty"`
}
