package model

import "gorm.io/gorm"

type Company struct {
	gorm.Model
	UserID      int
	User        User   `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ContactName string `gorm:"type:varchar(30);not null;" json:"contactName" validate:"required"`
	CompanyName string `gorm:"type:varchar(30);not null;" json:"companyName" validate:"required"`
	Country     string `gorm:"type:integer;not null;" json:"country" validate:"required"`
	Email       string `gorm:"type:varchar(35);not null;" json:"email" validate:"required"`
	CellPhone   string `gorm:"type:varchar(18);not null;" json:"cellPhone" validate:"required"`
	Phone       string `gorm:"type:varchar(18);not null;" json:"phone" validate:"required"`
}
