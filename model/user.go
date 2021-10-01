package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(20);not null;unique" json:"name,omitempty" validate:"required"`
	Password string `gorm:"size:255" json:"pass,omitempty" validate:"required"`
	Role     string `gorm:"type:varchar(20);not null;" json:"role,omitempty" validate:"required"`
}

type UserFormData struct {
	User
	Pass        string `json:"pass,omitempty" binding:"required"`
	ConfirmPass string `json:"confirmPass,omitempty" binding:"required"`
}

func (user *User) BeforeCreate(db *gorm.DB) (err error) {
	if pw, err := bcrypt.GenerateFromPassword([]byte(string(user.Password)), 14); err == nil {
		user.Password = string(pw)
	}

	return
}

func (user *User) BeforeUpdate(db *gorm.DB) (err error) {
	if pw, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14); err == nil {
		user.Password = string(pw)
	}

	return
}

func (user *User) FindUser(db *gorm.DB, username interface{}) error {
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return err
	} else {
		return err
	}
}
