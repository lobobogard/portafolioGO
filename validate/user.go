package validate

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/portafolioLP/model"
	"gorm.io/gorm"
)

var validate *validator.Validate

func ValidateUser(UserFormData *model.UserFormData, DB *gorm.DB, w http.ResponseWriter, r *http.Request) error {
	validate = validator.New()

	myUsername := UserFormData.Username

	var user model.User
	DB.Where("username = ?", myUsername).First(&user)
	if user.Username != "" {
		return errors.New("the user name exist")
	}

	err := validate.Var(myUsername, "required")
	if err != nil {
		return errors.New("the user name is required")
	}

	myPassword := UserFormData.Pass
	fmt.Println(myPassword)
	err = validate.Var(myPassword, "required")
	if err != nil {
		return errors.New("the password is required")
	}

	myConfirmPassword := UserFormData.ConfirmPass
	err = validate.Var(myConfirmPassword, "required")
	if err != nil {
		return errors.New("the confirm password is required")
	}

	if myPassword != myConfirmPassword {
		return errors.New("the confirm password and password are diff")
	}

	return err
}

func ValidateStruct(user *model.User, w http.ResponseWriter, r *http.Request) error {
	validate = validator.New()
	return validate.Struct(user)
}
