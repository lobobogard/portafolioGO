package validate

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/portafolioLP/model"
)

var validate *validator.Validate

func Validate(user model.User) {
	validate = validator.New()
	myUsername := user.Username
	errs := validate.Var(myUsername, "required")

	if errs != nil {
		fmt.Println(errs) // output: Key: "" Error:Field validation for "" failed on the "email" tag
	}
}

func ValidateStruct(user *model.User, w http.ResponseWriter, r *http.Request) error {
	validate = validator.New()
	return validate.Struct(user)
}
