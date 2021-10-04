package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/portafolioLP/model"
	"github.com/portafolioLP/validate"
	"gorm.io/gorm"
)

func CreatePerfil(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
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
		message, httpStatus := savePerfil(DB, w, r, PerfilFormData)
		respondJSON(w, httpStatus, message)
		// respondJSON(w, http.StatusAccepted, "the profile was created successfully")
	}
}

func savePerfil(DB *gorm.DB, w http.ResponseWriter, r *http.Request, PerfilFormData *model.ReqPerfil) (string, int) {

	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err.Error(), http.StatusInternalServerError
	}

	val, perfil := createPerfil(PerfilFormData, tx)
	if val != "" {
		return val, http.StatusInternalServerError
	}

	val = createServer(PerfilFormData, tx, DB, perfil)
	if val != "" {
		return val, http.StatusInternalServerError
	}

	val = createBackEnd(PerfilFormData, tx, DB, perfil)
	if val != "" {
		return val, http.StatusInternalServerError
	}

	val = createFrontEnd(PerfilFormData, tx, DB, perfil)
	if val != "" {
		return val, http.StatusInternalServerError
	}

	val, database := createDatabase(PerfilFormData, tx, perfil)
	if val != "" {
		return val, http.StatusInternalServerError
	}

	val = updatePerfilIDS(DB, tx, perfil, database)
	if val != "" {
		return val, http.StatusInternalServerError
	}

	tx.Commit()
	return "profile was created successfully", http.StatusOK
}

func createPerfil(PerfilFormData *model.ReqPerfil, tx *gorm.DB) (string, model.Perfil) {
	perfil := model.Perfil{}
	perfil.CompanyID = PerfilFormData.CompanyID
	perfil.SystemOperativeID = PerfilFormData.SystemOperativeID

	if err := tx.Create(&perfil).Error; err != nil {
		tx.Rollback()
		return "Error in system create perfil", perfil
	}
	fmt.Println("id perfil", perfil.ID)
	return "", perfil
}

func updatePerfilIDS(DB *gorm.DB, tx *gorm.DB, perfil model.Perfil, database model.DataBase) string {

	DB.First(&perfil)
	perfil.DataBaseID = &database.ID
	perfil.ServerdeployPerfilID = &perfil.ID
	perfil.BackEndID = &perfil.ID
	perfil.FrontEndID = &perfil.ID

	if perfil.ID == 0 {
		tx.Rollback()
		return "Error in system perfilID not exist"
	}

	if err := tx.Save(&perfil).Error; err != nil {
		tx.Rollback()
		fmt.Println(err)
		return "Error in system update perfil ID"
	}

	return ""
}

func createServer(PerfilFormData *model.ReqPerfil, tx *gorm.DB, DB *gorm.DB, perfil model.Perfil) string {
	var catServer []model.CatServer
	var server []model.Servers
	DB.Find(&catServer)
	for _, v := range catServer {
		var s model.Servers
		s.PerfilID = perfil.ID
		s.CatserverID = v.ID
		s.Status = ItemExists(PerfilFormData.Server, v.ID)
		server = append(server, s)
	}

	if err := tx.Create(&server).Error; err != nil {
		tx.Rollback()
		return "Error in system create server"
	}

	return ""
}

func createDatabase(PerfilFormData *model.ReqPerfil, tx *gorm.DB, perfil model.Perfil) (string, model.DataBase) {
	database := model.DataBase{PerfilID: perfil.ID,
		Mysql: PerfilFormData.Mysql, Mariadb: PerfilFormData.Mariadb, Postgresql: PerfilFormData.Postgresql,
		Mongodb: PerfilFormData.Mongodb, Redis: PerfilFormData.Redis, Sqlite: PerfilFormData.Sqlite,
	}

	if err := tx.Create(&database).Error; err != nil {
		tx.Rollback()
		return "Error in system create database", database
	}

	return "", database
}

func createFrontEnd(PerfilFormData *model.ReqPerfil, tx *gorm.DB, DB *gorm.DB, perfil model.Perfil) string {
	var catFrontEnd []model.CatFrontEnd
	var frontEnd []model.FrontEnd
	DB.Find(&catFrontEnd)
	for _, v := range catFrontEnd {
		var s model.FrontEnd
		s.PerfilID = perfil.ID
		s.CatfrontendID = v.ID
		s.Status = ItemExists(PerfilFormData.FrontEnd, v.ID)
		frontEnd = append(frontEnd, s)
	}

	if err := tx.Create(&frontEnd).Error; err != nil {
		tx.Rollback()
		return "Error in system create frontEnd"
	}

	return ""
}

func createBackEnd(PerfilFormData *model.ReqPerfil, tx *gorm.DB, DB *gorm.DB, perfil model.Perfil) string {
	var catBackEnd []model.CatBackEnd
	var backEnd []model.BackEnd
	DB.Find(&catBackEnd)
	for _, v := range catBackEnd {
		var s model.BackEnd
		s.PerfilID = perfil.ID
		s.CatbackendID = v.ID
		s.Status = ItemExists(PerfilFormData.BackEnd, v.ID)
		backEnd = append(backEnd, s)
	}

	if err := tx.Create(&backEnd).Error; err != nil {
		tx.Rollback()
		return "Error in system create backEnd"
	}

	return ""
}

func ItemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}

	return false
}
