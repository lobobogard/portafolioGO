package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/portafolioLP/model"
	"github.com/portafolioLP/validate"
	"gorm.io/gorm"
)

func CreateCompany(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	userJWT := DecodeSessionUserJWT(w, r)

	var user model.User
	DB.First(&user, "username = ?", userJWT.Username)

	CompanyFormData := &model.Company{}
	CompanyFormData.UserID = user.ID

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(CompanyFormData); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := validate.ValidateCompany(CompanyFormData, w, r); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := DB.Create(&CompanyFormData).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	flagSendEmail, dataEmail := Emails(DB, w, r)
	if flagSendEmail {
		respondJSON(w, http.StatusCreated, dataEmail)
	} else {
		respondJSON(w, http.StatusCreated, "Created company success")
	}
}

func UpdateCompany(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	CompanyFormData := &model.Company{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(CompanyFormData); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := validate.ValidateCompany(CompanyFormData, w, r); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	var company model.Company
	DB.First(&company, vars["companyID"])
	company.CompanyName = CompanyFormData.CompanyName
	company.ContactName = CompanyFormData.ContactName
	company.Country_id = CompanyFormData.Country_id
	company.Email = CompanyFormData.Email
	company.CellPhone = CompanyFormData.CellPhone
	company.Phone = CompanyFormData.Phone
	DB.Save(&company)
	respondJSON(w, http.StatusCreated, "Update company success")
}

func FindCompany(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	user := DecodeSessionUserDB(DB, w, r)
	companyName := r.URL.Query().Get("company")
	contactName := r.URL.Query().Get("contact")
	email := r.URL.Query().Get("email")

	company := &[]model.Company{}

	base := DB.Where("user_id = ?", user.ID)
	if companyName != "" && companyName != "null" {
		base = base.Where("company_name LIKE ?", "%"+companyName+"%")
	}

	if contactName != "" && contactName != "null" {
		base = base.Where("contact_name LIKE ?", "%"+contactName+"%")
	}

	if email != "" && email != "null" {
		base = base.Where("email LIKE ?", "%"+email+"%")
	}
	base.Joins("CatCountry").Find(&company)
	respondJSON(w, http.StatusCreated, company)
}

func DeleteCompany(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyID := vars["companyID"]
	err := deleteCompany(DB, companyID)

	if err != nil {
		respondJSON(w, http.StatusBadRequest, "Error system comunicate with admin")
	} else {
		respondJSON(w, http.StatusCreated, "Delete successfully")
	}
}

func deleteCompany(DB *gorm.DB, companyID string) error {
	company := &model.Company{}
	DB.Find(company, companyID)
	perfil := &model.Perfil{}
	DB.Find(perfil, "company_id", companyID)

	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Where("perfil_id = ?", perfil.ID).Unscoped().Delete(&model.Servers{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("perfil_id = ?", perfil.ID).Unscoped().Delete(&model.DataBase{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("perfil_id = ?", perfil.ID).Unscoped().Delete(&model.BackEnd{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("perfil_id = ?", perfil.ID).Unscoped().Delete(&model.FrontEnd{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("id = ?", perfil.ID).Unscoped().Delete(&perfil).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Unscoped().Delete(&company).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
