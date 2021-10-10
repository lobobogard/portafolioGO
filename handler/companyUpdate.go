package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/portafolioLP/model"
	"gorm.io/gorm"
)

type companyUpdateData struct {
	Company    []model.Company
	CatCountry []model.CatCountry
}

func GetCompanyUpdate(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var data companyUpdateData
	DB.Find(&data.Company, vars["companyID"])
	DB.Select("id", "country").Find(&data.CatCountry)
	respondJSON(w, http.StatusOK, data)
}
