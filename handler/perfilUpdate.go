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

type responsePerfil struct {
	Company            []model.Company            `json:"Company"`
	CatSystemOperative []model.CatSystemOperative `json:"CatSystemOperative"`
	CatServer          []model.CatServer          `json:"CatServer"`
	BackEnd            []model.CatBackEnd         `json:"BackEnd"`
	FrontEnd           []model.CatFrontEnd        `json:"FrontEnd"`
	GetDataPerfil      GetDataPerfil              `json:"getDataPerfil"`
}

type GetDataPerfil struct {
	Perfil                 model.Perfil             `json:"Perfil"`
	CatSystemOperativeData model.CatSystemOperative `json:"CatSystemOperativeData"`
	CatServerNames         []string                 `json:"CatServerNames"`
	BackEndIDS             []model.CatBackEnd       `json:"BackEndIDS"`
	FrontEndIDS            []model.CatFrontEnd      `json:"FrontEndIDS"`
	GetDataBase            GetDataBase
}

type GetDataBase struct {
	Mysql      bool `json:"Mysql"`
	Mariadb    bool `json:"Mariadb"`
	Postgresql bool `json:"Postgresql"`
	Mongodb    bool `json:"Mongodb"`
	Redis      bool `json:"Redis"`
	Sqlite     bool `json:"Sqlite"`
}

func MountPerfil(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	user := DecodeSessionUserDB(DB, w, r)
	vars := mux.Vars(r)
	perfilID := vars["perfilID"]
	var resp responsePerfil
	DB.Select("id", "company_name").Find(&resp.Company, "user_id", user.ID)
	DB.Find(&resp.BackEnd)
	DB.Find(&resp.FrontEnd)
	DB.Find(&resp.CatServer)
	DB.Find(&resp.CatSystemOperative)

	DB.Select("data_bases.*").
		Table("perfils").
		Joins("inner join data_bases on data_bases.perfil_id = perfils.data_base_id").
		Where("perfils.id = (?)", perfilID).
		Scan(&resp.GetDataPerfil.GetDataBase)

	DB.Select("company_id").First(&resp.GetDataPerfil.Perfil, perfilID)

	DB.Select("cat_system_operatives.id, cat_system_operatives.system_operative").
		Table("perfils").
		Joins("inner join cat_system_operatives on perfils.cat_system_operative_id = cat_system_operatives.id").
		Where("perfils.id = (?)", perfilID).
		Scan(&resp.GetDataPerfil.CatSystemOperativeData)

	DB.Select("cat_servers.server").
		Table("perfils").
		Joins("inner join servers on servers.perfil_id = perfils.serverdeploy_perfil_id").
		Joins("inner join cat_servers on cat_servers.id = servers.catserver_id").
		Where("perfils.id = (?) and servers.status in (?)", perfilID, 1).
		Scan(&resp.GetDataPerfil.CatServerNames)

	DB.Select("cat_back_ends.id, cat_back_ends.back_end").
		Table("perfils").
		Joins("inner join back_ends on back_ends.perfil_id = perfils.back_end_id ").
		Joins("inner join cat_back_ends on cat_back_ends.id = back_ends.catbackend_id").
		Where("perfils.id = (?) and back_ends.status in (?)", perfilID, 1).
		Scan(&resp.GetDataPerfil.BackEndIDS)

	DB.Select("cat_front_ends.id, cat_front_ends.front_end").
		Table("perfils").
		Joins("inner join front_ends on front_ends.perfil_id = perfils.front_end_id ").
		Joins("inner join cat_front_ends on cat_front_ends.id = front_ends.catfrontend_id").
		Where("perfils.id = (?) and front_ends.status in (?)", perfilID, 1).
		Scan(&resp.GetDataPerfil.FrontEndIDS)

	fmt.Println(resp.GetDataPerfil.CatSystemOperativeData)

	respondJSON(w, http.StatusOK, resp)
}

// trabajando update *********************************************************
func UpdatePerfil(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	PerfilFormData := &model.ReqPerfil{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(PerfilFormData); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := validate.ValidatePerfil(PerfilFormData, w, r); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
	} else {
		message, httpStatus := saveUpdatePerfil(DB, w, r, PerfilFormData, vars["perfilID"])
		respondJSON(w, httpStatus, message)
	}
}

func saveUpdatePerfil(DB *gorm.DB, w http.ResponseWriter, r *http.Request, PerfilFormData *model.ReqPerfil, perfilID string) (string, int) {

	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err.Error(), http.StatusInternalServerError
	}

	val, perfil := updatePerfil(PerfilFormData, tx, perfilID)
	if val != "" {
		return val, http.StatusInternalServerError
	}

	fmt.Println(perfil)

	val = updateServer(PerfilFormData, tx, DB, perfil)
	if val != "" {
		return val, http.StatusInternalServerError
	}

	val = updateDatabase(PerfilFormData, tx, perfil)
	if val != "" {
		return val, http.StatusInternalServerError
	}

	val = updateBackEnd(PerfilFormData, tx, DB, perfil)
	if val != "" {
		return val, http.StatusInternalServerError
	}

	val = updateFrontEnd(PerfilFormData, tx, DB, perfil)
	if val != "" {
		return val, http.StatusInternalServerError
	}

	tx.Commit()
	return "perfil updated successfully", http.StatusOK
}

func updatePerfil(PerfilFormData *model.ReqPerfil, tx *gorm.DB, perfilID string) (string, model.Perfil) {

	perfil := model.Perfil{}
	tx.First(&perfil, perfilID)
	perfil.CompanyID = PerfilFormData.CompanyID
	perfil.CatSystemOperativeID = PerfilFormData.SystemOperativeID

	if err := tx.Save(&perfil).Error; err != nil {
		tx.Rollback()
		return "Error in system create perfil", perfil
	}
	fmt.Println("id perfil", perfil.ID)
	return "", perfil
}

func updateServer(PerfilFormData *model.ReqPerfil, tx *gorm.DB, DB *gorm.DB, perfil model.Perfil) string {
	var catServer []model.CatServer
	var server []model.Servers
	DB.Find(&catServer)
	for _, v := range catServer {
		var s model.Servers
		tx.Where("perfil_id = (?) and catserver_id = (?)", perfil.ID, v.ID).First(&s)
		s.Status = ItemExists(PerfilFormData.Server, v.ID)
		server = append(server, s)
	}

	if err := tx.Save(&server).Error; err != nil {
		tx.Rollback()
		return "Error in system create server"
	}

	return ""
}

func updateDatabase(PerfilFormData *model.ReqPerfil, tx *gorm.DB, perfil model.Perfil) string {
	var database model.DataBase
	tx.Where("perfil_id = (?)", perfil.ID).First(&database)
	database.Mysql = PerfilFormData.Mysql
	database.Mariadb = PerfilFormData.Mariadb
	database.Postgresql = PerfilFormData.Postgresql
	database.Mongodb = PerfilFormData.Mongodb
	database.Redis = PerfilFormData.Redis
	database.Sqlite = PerfilFormData.Sqlite

	if err := tx.Save(&database).Error; err != nil {
		tx.Rollback()
		return "Error in system create database"
	}

	return ""
}

func updateBackEnd(PerfilFormData *model.ReqPerfil, tx *gorm.DB, DB *gorm.DB, perfil model.Perfil) string {
	var catBackEnd []model.CatBackEnd
	var backEnd []model.BackEnd
	DB.Find(&catBackEnd)
	for _, v := range catBackEnd {
		var s model.BackEnd
		tx.Where("perfil_id = (?) and catbackend_id = (?)", perfil.ID, v.ID).First(&s)
		s.Status = ItemExists(PerfilFormData.BackEnd, v.ID)
		backEnd = append(backEnd, s)
	}

	if err := tx.Save(&backEnd).Error; err != nil {
		tx.Rollback()
		return "Error in system create backEnd"
	}

	return ""
}

func updateFrontEnd(PerfilFormData *model.ReqPerfil, tx *gorm.DB, DB *gorm.DB, perfil model.Perfil) string {
	var catFrontEnd []model.CatFrontEnd
	var frontEnd []model.FrontEnd
	DB.Find(&catFrontEnd)
	for _, v := range catFrontEnd {
		var s model.FrontEnd
		tx.Where("perfil_id = (?) and catfrontend_id = (?)", perfil.ID, v.ID).First(&s)
		s.Status = ItemExists(PerfilFormData.FrontEnd, v.ID)
		frontEnd = append(frontEnd, s)
	}

	if err := tx.Save(&frontEnd).Error; err != nil {
		tx.Rollback()
		return "Error in system create frontEnd"
	}

	return ""
}
