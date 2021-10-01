package handler

import (
	"encoding/json"
	"net/http"

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

	if err := validate.ValidatePerfil(CompanyFormData, w, r); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := DB.Create(&CompanyFormData).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, "Created company success")
}
