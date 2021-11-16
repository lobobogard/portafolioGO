package handler

import (
	"encoding/json"
	"net/http"

	"github.com/portafolioLP/model"
	"gorm.io/gorm"
)

func GetConcurrency(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	userJWT := DecodeSessionUserNotVerificateJWT(w, r)
	var concurrency model.ConfConcurrency
	DB.Where("username", userJWT.Id).First(&concurrency)
	respondJSON(w, http.StatusCreated, concurrency)
}

func Email(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	userJWT := DecodeSessionUserNotVerificateJWT(w, r)
	var ConcurrencyFormData model.ConfConcurrency

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&ConcurrencyFormData); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()
	var concurrency model.ConfConcurrency
	if err := DB.Where("username", userJWT.Id).First(&concurrency).Error; err != nil {
		ConcurrencyFormData.Username = userJWT.Id
		DB.Create(&ConcurrencyFormData)
		respondJSON(w, http.StatusCreated, "changed send email successfully")
	} else {
		concurrency.Username = userJWT.Id
		concurrency.SendEmail = ConcurrencyFormData.SendEmail
		DB.Where("username", userJWT.Id).Save(&concurrency)
		respondJSON(w, http.StatusCreated, "changed send email successfully")
	}

}

func Concurrency(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	userJWT := DecodeSessionUserNotVerificateJWT(w, r)
	var ConcurrencyFormData model.ConfConcurrency

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&ConcurrencyFormData); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()
	var concurrency model.ConfConcurrency
	if err := DB.Where("username", userJWT.Id).First(&concurrency).Error; err != nil {
		ConcurrencyFormData.Username = userJWT.Id
		DB.Create(&ConcurrencyFormData)
		respondJSON(w, http.StatusCreated, "changed concurrency successfully")
	} else {
		concurrency.Username = userJWT.Id
		concurrency.Concurrency = ConcurrencyFormData.Concurrency
		DB.Where("username", userJWT.Id).Save(&concurrency)
		respondJSON(w, http.StatusCreated, "changed concurrency successfully")
	}

}
