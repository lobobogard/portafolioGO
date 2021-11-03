package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/portafolioLP/model"
	"gorm.io/gorm"
)

type PerfilFindData struct {
	Company            []model.Company
	CatSystemOperative []model.CatSystemOperative
	CatBackEnd         []model.CatBackEnd
}

func FindMountedPerfil(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var data PerfilFindData
	DB.Find(&data.Company)
	DB.Find(&data.CatSystemOperative)
	DB.Find(&data.CatBackEnd)
	respondJSON(w, http.StatusCreated, data)
}

type Result struct {
	ID                      int
	Company_id              string
	Cat_system_operative_id int
	Company_name            string
	System_operative        string
	Backend                 string
}

func FindPerfil(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var result []Result
	user := DecodeSessionUserDB(DB, w, r)
	companyID := r.URL.Query().Get("companyID")
	sysOperativeID := r.URL.Query().Get("sysOperativeID")
	backEnd := r.URL.Query().Get("backEnd")

	base := DB.Where("user_id = ? and back_ends.status = ?", user.ID, "1")
	if companyID != "" && companyID != "null" {
		base = base.Where("company_id", companyID)
	}

	if sysOperativeID != "" && sysOperativeID != "null" {
		base = base.Where("cat_system_operative_id", sysOperativeID)
	}

	if backEnd != "" && backEnd != "null" {
		backEndArray := convertStringToArrayInt(backEnd)
		base = base.Where("back_ends.catbackend_id in (?)", backEndArray)
	}

	base.Table("perfils").
		Select("perfils.id", "company_id", "companies.company_name", "cat_system_operatives.system_operative, GROUP_CONCAT(cat_back_ends.back_end) as backend").
		Joins("inner join back_ends on back_ends.perfil_id = perfils.back_end_id ").
		Joins("inner join cat_back_ends on cat_back_ends.id = back_ends.catbackend_id").
		Joins("inner join companies on companies.id = perfils.company_id").
		Joins("inner join cat_system_operatives on cat_system_operatives.id = perfils.cat_system_operative_id").
		Group("perfils.id,companies.company_name,cat_system_operatives.system_operative").
		Scan(&result)

	respondJSON(w, http.StatusCreated, result)
}

func DeletePerfil(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	perfilID := vars["perfilID"]
	err := deletePerfil(DB, perfilID)

	if err != nil {
		respondJSON(w, http.StatusBadRequest, "Error system comunicate with admin")
	} else {
		respondJSON(w, http.StatusCreated, "Delete successfully")
	}
}

func convertStringToArrayInt(backEnd string) []int {
	backEnd = string("[") + backEnd + string("]")
	var backEndArray []int
	err := json.Unmarshal([]byte(backEnd), &backEndArray)
	if err != nil {
		log.Fatal(err)
	}
	return backEndArray
}

func deletePerfil(DB *gorm.DB, perfilID string) error {
	perfil := &model.Perfil{}
	DB.First(perfil, perfilID)

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

	if err := tx.Unscoped().Delete(perfil).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
