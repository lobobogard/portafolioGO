package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/portafolioLP/authentication"
	"github.com/portafolioLP/model"
	"gorm.io/gorm"
)

type TokenRefresh struct {
	Token string
}

func ValidateTokenRefresh(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("tokenRefresh")
	if err != nil {
		jsonStatusUnauthorized("Not found cookie", w)
		return
	}
	tokenStr := cookie.Value
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return authentication.GetPublicKey(), nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)
	username := fmt.Sprintf("%v", claims["username"])

	if ok && token.Valid {
		var user model.User
		if user.FindUser(db, claims["username"]) != nil {
			jsonStatusUnauthorized("User not Exist", w)
			return
		}
	}

	if model.ExistUserRedisToken(username, tokenStr) {
		jsonStatusUnauthorized("TokenRefresh Unauthorized", w)
		return
	}

	if err != nil {
		switch err.(type) {
		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				jsonStatusUnauthorized("Your tokenRefresh has expired", w)
				return
			case jwt.ValidationErrorSignatureInvalid:
				jsonStatusUnauthorized("TokenRefresh signature does not match", w)
				return
			default:
				jsonStatusUnauthorized("Your tokenRefresh is not valid", w)
				return
			}
		default:
			jsonStatusUnauthorized("Your tokenRefresh is not valid", w)
			return
		}
	}
	if token.Valid {
		regenerateToken(db, username, w)
	} else {
		jsonStatusUnauthorized("Unauthorized", w)
	}
}

func regenerateToken(db *gorm.DB, username string, w http.ResponseWriter) {
	var user, userDB model.User

	db.Where("username = ?", username).First(&userDB)
	user.Password = ""
	user.Role = userDB.Role
	user.Username = userDB.Username
	token := authentication.GenerateJWT(user)
	tokenRefresh := authentication.GenerateRefreshJWT(w, userDB)
	model.RedisRefreshToken(user, tokenRefresh)
	result := model.ResponseToken{Token: token}
	jsonResult, err := json.Marshal(result)
	if err != nil {
		fmt.Fprintf(w, "Error al generar el json")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResult)
}

func DeleteTokenRefreshRedis(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("tokenRefresh")
	if err != nil {
		jsonStatusUnauthorized("Not found cookie", w)
		return
	}
	tokenStr := cookie.Value
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return authentication.GetPublicKey(), nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)
	username := fmt.Sprintf("%v", claims["username"])

	if ok && token.Valid {
		var user model.User
		if user.FindUser(db, claims["username"]) != nil {
			jsonStatusUnauthorized("User not Exist", w)
			return
		}
	}

	if model.ExistUserRedisToken(username, tokenStr) {
		jsonStatusUnauthorized("TokenRefresh Unauthorized", w)
		return
	}

	if err != nil {
		switch err.(type) {
		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				jsonStatusUnauthorized("Your tokenRefresh has expired", w)
				return
			case jwt.ValidationErrorSignatureInvalid:
				jsonStatusUnauthorized("TokenRefresh signature does not match", w)
				return
			default:
				jsonStatusUnauthorized("Your tokenRefresh is not valid", w)
				return
			}
		default:
			jsonStatusUnauthorized("Your tokenRefresh is not valid", w)
			return
		}
	}
	if token.Valid {
		model.RedisDeleteRefreshToken(username)
		jsonStatusAccepted("Your tokenRefresh was deleted successfully", w)
	} else {
		jsonStatusUnauthorized("Unauthorized", w)
	}
}

func jsonStatusUnauthorized(result string, w http.ResponseWriter) {
	jsonResult, err := json.Marshal(result)
	if err != nil {
		fmt.Fprintf(w, "Error al generar el json")
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResult)
}

func jsonStatusAccepted(result string, w http.ResponseWriter) {
	jsonResult, err := json.Marshal(result)
	if err != nil {
		fmt.Fprintf(w, "Error al generar el json")
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResult)
}
