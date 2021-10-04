package validate

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/portafolioLP/model"
)

func ValidatePerfil(PerfilFormData *model.ReqPerfil, w http.ResponseWriter, r *http.Request) error {
	validate = validator.New()

	CompanyID := PerfilFormData.CompanyID
	err := validate.Var(CompanyID, "required")
	if err != nil {
		return errors.New("the select the company is required")
	}

	SystemOperativeID := PerfilFormData.SystemOperativeID
	err = validate.Var(SystemOperativeID, "required")
	if err != nil {
		return errors.New("the select the system operative is required")
	}

	Server := PerfilFormData.Server
	if len(Server) == 0 {
		fmt.Println("cantidad", len(Server))
		return errors.New("select one server as minimum")
	}

	if !PerfilFormData.Mariadb && !PerfilFormData.Mongodb && !PerfilFormData.Mysql && !PerfilFormData.Postgresql && !PerfilFormData.Redis && !PerfilFormData.Sqlite {
		return errors.New("select one database as minimum")
	}

	BackEnd := PerfilFormData.BackEnd
	if len(BackEnd) == 0 {
		fmt.Println("cantidad", len(BackEnd))
		return errors.New("select one backend as minimum")
	}

	FrontEnd := PerfilFormData.FrontEnd
	if len(FrontEnd) == 0 {
		return errors.New("select one frontend as minimum")
	}

	return err

}
