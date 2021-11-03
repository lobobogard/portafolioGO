package model

import "gorm.io/gorm"

type CatFrontEnd struct {
	gorm.Model
	FrontEnd string `gorm:"type:varchar(20);not null;unique" json:"frontEnd" validate:"required"`
}
