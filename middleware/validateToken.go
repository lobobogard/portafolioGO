package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/request"
	"github.com/portafolioLP/authentication"
	"github.com/portafolioLP/model"
)

func ValidateToken(w http.ResponseWriter, r *http.Request) (bool, string) {
	var mensaje string
	token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &model.Claim{}, func(token *jwt.Token) (interface{}, error) {
		return authentication.GetPublicKey(), nil
	})

	if err != nil {
		switch err.(type) {
		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				mensaje = "Your token has expired"
				return false, mensaje
			case jwt.ValidationErrorSignatureInvalid:
				mensaje = "Token signature does not match"
				return false, mensaje
			default:
				mensaje = "Your token is not valid"
				return false, mensaje
			}
		default:
			mensaje = "Your token is not valid"
			return false, mensaje
		}
	}
	if token.Valid {
		mensaje = "correct"
		return true, mensaje
	} else {
		mensaje = "Unauthorized"
		return false, mensaje
	}
}

type Data struct {
	Mensaje  string
	Validate bool
}

func Logging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var datos Data
		datos.Validate, datos.Mensaje = ValidateToken(w, r)
		if datos.Validate {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(datos)
		}
	}
}
