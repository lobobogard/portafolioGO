package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/portafolioLP/model"
	"gorm.io/gorm"
)

type ReqPerfil struct {
	CompanyID         uint
	SystemOperativeID uint
	Server            []uint
	Mysql             bool
	Mariadb           bool
	Postgresql        bool
	Mongodb           bool
	Redis             bool
	Sqlite            bool
}

func CreatePerfil(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	PerfilFormData := &ReqPerfil{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(PerfilFormData); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	message, httpStatus := savePerfil(DB, w, r, PerfilFormData)

	respondJSON(w, httpStatus, message)
	// respondJSON(w, http.StatusAccepted, "correcto")
}

func savePerfil(DB *gorm.DB, w http.ResponseWriter, r *http.Request, PerfilFormData *ReqPerfil) (string, int) {

	// if err := validate.ValidatePerfil(CompanyFormData, w, r); err != nil {
	// 	respondError(w, http.StatusBadRequest, err.Error())
	// 	return
	// }

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

	val, database := createDatabase(PerfilFormData, tx, perfil)
	if val != "" {
		return val, http.StatusInternalServerError
	}

	val = updatePerfilIDS(DB, tx, perfil, database)
	if val != "" {
		return val, http.StatusInternalServerError
	}

	tx.Commit()
	return "Create perfil succeful", http.StatusOK
}

func createPerfil(PerfilFormData *ReqPerfil, tx *gorm.DB) (string, model.Perfil) {
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

func createDatabase(PerfilFormData *ReqPerfil, tx *gorm.DB, perfil model.Perfil) (string, model.DataBase) {
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

func createServer(PerfilFormData *ReqPerfil, tx *gorm.DB, DB *gorm.DB, perfil model.Perfil) string {
	var catServer []model.CatServer
	var server []model.Servers
	DB.Find(&catServer)
	for _, v := range catServer {
		var s model.Servers
		s.PerfilID = perfil.ID
		s.CatserverID = v.ID
		s.StatusServer = ItemExists(PerfilFormData.Server, v.ID)
		server = append(server, s)
	}

	if err := tx.Create(&server).Error; err != nil {
		tx.Rollback()
		return "Error in system create server"
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

func updatePerfilIDS(DB *gorm.DB, tx *gorm.DB, perfil model.Perfil, database model.DataBase) string {

	DB.First(&perfil)
	perfil.DataBaseID = &database.ID
	perfil.ServerdeployPerfilID = &perfil.ID

	if perfil.ID == 0 {
		tx.Rollback()
		return "Error in system perfilID not exist"
	}

	if err := tx.Save(&perfil).Error; err != nil {
		tx.Rollback()
		return "Error in system update perfil databaseID"
	}

	return ""
}
