package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Reading body fail")
	}

	tokenBody := TokenRefresh{}
	err = json.Unmarshal(body, &tokenBody)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return
	}

	tokenStr := tokenBody.Token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return authentication.GetPublicKey(), nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		var user model.User
		if user.FindUser(db, claims["username"]) != nil {
			fmt.Fprintf(w, "User not Found")
			return
		}
	}

	username := fmt.Sprintf("%v", claims["username"])
	fmt.Println(model.ExistUserRedisToken(username, tokenBody.Token))
	if model.ExistUserRedisToken(username, tokenBody.Token) {
		fmt.Fprintf(w, "Su tokenRefresh Unauthorized")
		return
	}

	if err != nil {
		switch err.(type) {
		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				fmt.Fprintf(w, "Su tokenRefresh ha expirado")
				return
			case jwt.ValidationErrorSignatureInvalid:
				fmt.Fprintf(w, "La firma del tokenRefresh no coincide")
				return
			default:
				fmt.Fprintf(w, "Su tokenRefresh no es valido")
				return
			}
		default:
			fmt.Fprintf(w, "Su tokenRefresh no es valido")
			return
		}
	}
	if token.Valid {
		regenerateToken(db, username, w)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorized")
	}
}

func regenerateToken(db *gorm.DB, username string, w http.ResponseWriter) {
	var user, userDB model.User

	db.Where("username = ?", username).First(&userDB)
	user.Password = ""
	user.Role = userDB.Role
	token := authentication.GenerateJWT(user, 5)
	tokenRefresh := authentication.GenerateRefreshJWT(user, 15)
	model.RedisRefreshToken(user, tokenRefresh)
	result := model.ResponseToken{Token: token, TokenRefresh: tokenRefresh}
	jsonResult, err := json.Marshal(result)
	if err != nil {
		fmt.Fprintf(w, "Error al generar el json")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResult)
}
