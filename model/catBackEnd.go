package model

import "gorm.io/gorm"

type CatBackEnd struct {
	gorm.Model
	BackEnd string `gorm:"type:varchar(20);not null;unique" json:"backEnd" validate:"required"`
}
