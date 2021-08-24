package model

import "gorm.io/gorm"

type Perfil struct {
	gorm.Model
	CompanyID       int
	Company         Company `gorm:"foreignKey:CompanyID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Skills          string
	LenguajeSpeak   string
	LenguajeProgram string
	DB              string
	Area            int //frontend backend devops
}
