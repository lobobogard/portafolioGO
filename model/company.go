package model

import "gorm.io/gorm"

type Company struct {
	gorm.Model
	UserID      int
	User        User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	NameRH      string
	CompanyName string
}
