package model

import "gorm.io/gorm"

type CatSystemOperative struct {
	gorm.Model
	SystemOperative string `gorm:"type:varchar(20);not null;unique" json:"systemOperative"`
}
