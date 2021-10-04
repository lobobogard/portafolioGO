package model

import "gorm.io/gorm"

type Perfil struct {
	gorm.Model
	CompanyID            uint    `gorm:"null" json:"companyID"`
	Company              Company `gorm:"foreignKey:CompanyID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	SystemOperativeID    uint    `gorm:"not null" json:"systemOperativeID"`
	ServerdeployPerfilID *uint   `gorm:"null" json:"ServerdeployPerfilID"`
	DataBaseID           *uint   `gorm:"null" json:"dataBaseID"`
	BackEndID            *uint   `gorm:"null" json:"backEndID"`
	FrontEndID           *uint   `gorm:"null" json:"frontEndID"`
}
