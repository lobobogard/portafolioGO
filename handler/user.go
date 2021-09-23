package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/portafolioLP/model"
	"github.com/portafolioLP/validate"
	"gorm.io/gorm"
)

func User(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	UserFormData := &model.UserFormData{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(UserFormData); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := validate.ValidateUser(UserFormData, w, r); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	UserFormData.User.Password = UserFormData.ConfirmPass
	if err := DB.Create(&UserFormData.User).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, UserFormData.User.Username)
}

func GetAllUsers(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	users := []model.User{}
	db.Find(&users)
	respondJSON(w, http.StatusOK, users)
}

func GetUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	user := getUserOr404(db, name, w, r)
	if user == nil {
		return
	}
	respondJSON(w, http.StatusOK, user)
}

func UpdateUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]
	user := getUserOr404(db, name, w, r)
	if user == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&user).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, user)
}

func DeleteUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	user := getUserOr404(db, name, w, r)
	if user == nil {
		return
	}
	fmt.Println(name, user)
	if err := db.Unscoped().Delete(&user).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}

func getUserOr404(db *gorm.DB, name string, w http.ResponseWriter, r *http.Request) *model.User {
	user := model.User{}
	if err := db.First(&user, model.User{Username: name}).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &user
}
