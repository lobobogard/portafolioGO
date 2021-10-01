package model

import "gorm.io/gorm"

type Servers struct {
	gorm.Model
	PerfilID     uint `gorm:"not null;" json:"PerfilID,omitempty"`
	CatserverID  uint `gorm:"not null;" json:"catserverID,omitempty"`
	StatusServer bool `gorm:"not null;" json:"statusServer,omitempty"`
}
