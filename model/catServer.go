package model

import "gorm.io/gorm"

type CatServer struct {
	gorm.Model
	Server string `gorm:"type:varchar(20);not null;unique" json:"server"`
}
