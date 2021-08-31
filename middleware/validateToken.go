package middleware

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/request"
	"github.com/portafolioLP/authentication"
	"github.com/portafolioLP/model"
)

func ValidateToken(w http.ResponseWriter, r *http.Request) bool {
	token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &model.Claim{}, func(token *jwt.Token) (interface{}, error) {
		return authentication.GetPublicKey(), nil
	})

	if err != nil {
		switch err.(type) {
		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				fmt.Fprintf(w, "Su token ha expirado")
				return false
			case jwt.ValidationErrorSignatureInvalid:
				fmt.Fprintf(w, "La firma del token no coincide")
				return false
			default:
				fmt.Fprintf(w, "Su token no es valido")
				return false
			}
		default:
			fmt.Fprintf(w, "Su token no es valido")
			return false
		}
	}
	if token.Valid {
		w.WriteHeader(http.StatusAccepted)
		// fmt.Fprintf(w, "Accepted")
		return true
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorized")
		return false
	}
}
