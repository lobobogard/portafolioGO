package model

import "gorm.io/gorm"

type ConfConcurrency struct {
	gorm.Model
	Username    string `gorm:"not null;" json:"userName,omitempty"`
	SendEmail   bool   `gorm:"not null;" json:"sendEmail"`
	Concurrency bool   `gorm:"not null;" json:"concurrency"`
}
