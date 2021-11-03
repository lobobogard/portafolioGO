package model

import "gorm.io/gorm"

type CatCountry struct {
	gorm.Model
	Country string `gorm:"type:varchar(20);not null;unique" json:"country" validate:"required"`
}
