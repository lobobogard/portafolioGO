package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/portafolioLP/authentication"
	"github.com/portafolioLP/model"
	"gorm.io/gorm"
)

func DecodeSessionUserJWT(w http.ResponseWriter, r *http.Request) model.Claim {
	tokenSession := r.Header.Get("Authorization")
	tokenStr := strings.Split(tokenSession, "Bearer ")[1]

	claims := jwt.MapClaims{}
	_, ok := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return authentication.GetPublicKey(), nil
	})

	if ok != nil {
		respondError(w, http.StatusBadRequest, "incorrect token session")
	}

	var tokenMap = make(map[string]interface{})
	for key, val := range claims {
		tokenMap[key] = val
	}

	jsonString, err := json.Marshal(tokenMap)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Error encoding user name")
	}

	user := model.Claim{}
	json.Unmarshal(jsonString, &user)
	return user
}

func DecodeSessionUserDB(DB *gorm.DB, w http.ResponseWriter, r *http.Request) model.User {
	tokenSession := r.Header.Get("Authorization")
	tokenStr := strings.Split(tokenSession, "Bearer ")[1]

	claims := jwt.MapClaims{}
	_, ok := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return authentication.GetPublicKey(), nil
	})

	if ok != nil {
		respondError(w, http.StatusBadRequest, "incorrect token session")
	}

	var tokenMap = make(map[string]interface{})
	for key, val := range claims {
		tokenMap[key] = val
	}

	jsonString, err := json.Marshal(tokenMap)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Error encoding user name")
	}

	userJWT := model.Claim{}
	json.Unmarshal(jsonString, &userJWT)

	var user model.User
	DB.First(&user, "username = ?", userJWT.Username)
	return user
}
