package validate

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/portafolioLP/model"
)

func ValidatePerfil(UserFormData *model.Company, w http.ResponseWriter, r *http.Request) error {
	validate = validator.New()

	CompanyName := UserFormData.CompanyName
	err := validate.Var(CompanyName, "required")
	if err != nil {
		return errors.New("the company name is required")
	}

	err = validate.Var(CompanyName, "max=20")
	if err != nil {
		return errors.New("the company name can't pass 20 characters")
	}

	ContactName := UserFormData.ContactName
	err = validate.Var(ContactName, "required")
	if err != nil {
		return errors.New("the contact name is required")
	}

	err = validate.Var(ContactName, "max=20")
	if err != nil {
		return errors.New("the contact name can't pass 20 characters")
	}

	Country := UserFormData.Country
	err = validate.Var(Country, "required")
	if err != nil {
		return errors.New("the country is required")
	}

	err = validate.Var(Country, "numeric")
	if err != nil {
		return errors.New("the country must be numeric")
	}

	Email := UserFormData.Email
	err = validate.Var(Email, "required")
	if err != nil {
		return errors.New("the email is required")
	}

	err = validate.Var(Email, "email")
	if err != nil {
		return errors.New("the email must be valid")
	}

	CellPhone := UserFormData.CellPhone
	err = validate.Var(CellPhone, "required")
	if err != nil {
		return errors.New("the cell phone is required")
	}

	err = validate.Var(CellPhone, "numeric")
	if err != nil {
		return errors.New("the cell phone must be numeric")
	}

	err = validate.Var(CellPhone, "max=18")
	if err != nil {
		return errors.New("the cell phone can't pass 18 characters")
	}

	Phone := UserFormData.Phone
	err = validate.Var(Phone, "required")
	if err != nil {
		return errors.New("the phone is required")
	}

	err = validate.Var(Phone, "numeric")
	if err != nil {
		return errors.New("the phone must be numeric")
	}

	err = validate.Var(Phone, "max=18")
	if err != nil {
		return errors.New("the phone can't pass 18 characters")
	}

	return err
}
