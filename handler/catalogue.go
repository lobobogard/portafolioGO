package handler

import (
	"net/http"

	"github.com/portafolioLP/model"
	"gorm.io/gorm"
)

func CataloguePerfil(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var country []model.CatCountry
	DB.Select("id", "country").Find(&country)

	respondJSON(w, http.StatusOK, country)
}
