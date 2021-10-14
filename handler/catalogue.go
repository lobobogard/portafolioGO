package handler

import (
	"net/http"

	"github.com/portafolioLP/model"
	"gorm.io/gorm"
)

type response struct {
	Company            []model.Company            `json:"Company"`
	CatSystemOperative []model.CatSystemOperative `json:"CatSystemOperative"`
	CatServer          []model.CatServer          `json:"CatServer"`
	BackEnd            []model.CatBackEnd         `json:"BackEnd"`
	FrontEnd           []model.CatFrontEnd        `json:"FrontEnd"`
}

type responseCompanyCreate struct {
	ConfConcurrency model.ConfConcurrency
	Country         []model.CatCountry
}

func CatalogueCountry(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var response responseCompanyCreate
	userJWT := DecodeSessionUserJWT(w, r)
	DB.Where("username", userJWT.Id).First(&response.ConfConcurrency)
	DB.Select("id", "country").Find(&response.Country)

	respondJSON(w, http.StatusOK, response)
}

func CatalogueCompany(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	user := DecodeSessionUserDB(DB, w, r)
	var resp response
	DB.Select("id", "company_name").Find(&resp.Company, "user_id", user.ID)
	DB.Find(&resp.CatServer)
	DB.Find(&resp.BackEnd)
	DB.Find(&resp.FrontEnd)
	DB.Find(&resp.CatSystemOperative)

	respondJSON(w, http.StatusOK, resp)
}
