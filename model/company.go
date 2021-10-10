package model

import "gorm.io/gorm"

type Company struct {
	gorm.Model
	UserID      uint
	User        User       `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ContactName string     `gorm:"type:varchar(30);not null;" json:"contactName,omitempty" validate:"required"`
	CompanyName string     `gorm:"type:varchar(30);not null;" json:"companyName,omitempty" validate:"required"`
	Country_id  uint       `gorm:"not null;" json:"Country_id,omitempty" validate:"required"`
	CatCountry  CatCountry `gorm:"foreignKey:Country_id"`
	Email       string     `gorm:"type:varchar(35);not null;" json:"email,omitempty" validate:"required"`
	CellPhone   string     `gorm:"type:varchar(18);not null;" json:"cellPhone,omitempty" validate:"required"`
	Phone       string     `gorm:"type:varchar(18);not null;" json:"phone,omitempty" validate:"required"`
}
