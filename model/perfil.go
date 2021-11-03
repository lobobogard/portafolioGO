package model

import "gorm.io/gorm"

type Perfil struct {
	gorm.Model
	CompanyID            uint               `gorm:"null" json:"companyID"`
	Company              Company            `gorm:"foreignKey:CompanyID"`
	CatSystemOperativeID uint               `gorm:"not null" json:"CatSystemOperativeID"`
	CatSystemOperative   CatSystemOperative `gorm:"foreignKey:CatSystemOperativeID"`
	ServerdeployPerfilID *uint              `gorm:"null" json:"ServerdeployPerfilID"`
	DataBaseID           *uint              `gorm:"null" json:"dataBaseID"`
	BackEndID            *uint              `gorm:"null" json:"backEndID"`
	FrontEndID           *uint              `gorm:"null" json:"frontEndID"`
}
