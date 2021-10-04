package model

import "gorm.io/gorm"

type DataBase struct {
	gorm.Model
	PerfilID   uint   `gorm:"not null;" json:"PerfilID,omitempty"`
	Perfil     Perfil `gorm:"foreignKey:PerfilID;"`
	Mysql      bool   `gorm:"not null;" json:"mysql"`
	Mariadb    bool   `gorm:"not null;" json:"mariadb"`
	Postgresql bool   `gorm:"not null;" json:"postgresql"`
	Mongodb    bool   `gorm:"not null;" json:"mongodb"`
	Redis      bool   `gorm:"not null;" json:"redis"`
	Sqlite     bool   `gorm:"not null;" json:"sqlite"`
}
