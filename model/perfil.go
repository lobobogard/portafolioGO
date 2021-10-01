package model

import "gorm.io/gorm"

type Perfil struct {
	gorm.Model
	CompanyID            uint               `gorm:"null" json:"companyID"`
	Company              Company            `gorm:"foreignKey:CompanyID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	SystemOperativeID    uint               `gorm:"not null" json:"systemOperativeID"`
	CatSystemOperative   CatSystemOperative `gorm:"foreignKey:SystemOperativeID;"`
	ServerdeployPerfilID *uint              `gorm:"null" json:"ServerdeployPerfilID"`
	Servers              Servers            `gorm:"foreignKey:ServerdeployPerfilID;"`
	DataBaseID           *uint              `gorm:"null" json:"dataBaseID"`
	DataBase             DataBase           `gorm:"foreignKey:DataBaseID;"`
	BackEndID            *uint              `gorm:"null" json:"backEndID"`
	BackEnd              BackEnd            `gorm:"foreignKey:BackEndID;"`
	FrontEndID           *uint              `gorm:"null" json:"frontEndID"`
	FrontEnd             FrontEnd           `gorm:"foreignKey:FrontEndID;"`
}
