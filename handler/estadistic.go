package handler

import (
	"fmt"
	"net/http"

	"github.com/portafolioLP/model"
	"gorm.io/gorm"
)

type DataEstadistic struct {
	Company []model.Company
	BackEnd []model.CatBackEnd
}

func MountEstadistic(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	user := DecodeSessionUserDB(DB, w, r)
	var mounted DataEstadistic
	DB.Select("companies.id, company_name").
		Table("companies").
		Joins("inner join users on  companies.user_id = users.id").
		Where("users.id = (?)", user.ID).
		Find(&mounted.Company)

	DB.Select("id, back_end").Find(&mounted.BackEnd)
	respondJSON(w, http.StatusAccepted, mounted)
}

type DataBackEnd struct {
	Catbackend_id int
	Back_end      string
	ReportBackEnd int
}

func Estadistic(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	user := DecodeSessionUserDB(DB, w, r)
	company := r.URL.Query().Get("company")
	backEnd := r.URL.Query().Get("backEnd")
	var result []DataBackEnd

	companyArray := convertStringToArrayInt(company)
	backEndArray := convertStringToArrayInt(backEnd)
	fmt.Println(companyArray)
	DB.Select("back_ends.catbackend_id,cat_back_ends.back_end, count(*) as ReportBackEnd").
		Table("back_ends").
		Joins("inner join perfils on perfils.id = back_ends.perfil_id").
		Joins("inner join companies on  companies.id = perfils.company_id").
		Joins("inner join cat_back_ends on back_ends.catbackend_id = cat_back_ends.id ").
		Where("companies.user_id = (?) and status = (?) and catbackend_id in (?) and company_id in (?)", user.ID, 1, backEndArray, companyArray).
		Group("back_ends.catbackend_id, cat_back_ends.back_end").
		Scan(&result)

	respondJSON(w, http.StatusAccepted, result)
}
