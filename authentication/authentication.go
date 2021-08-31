package authentication

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/request"
	"github.com/portafolioLP/model"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//crear llave privadas y publicas minuto 26
var (
	privateKey *rsa.PrivateKey // openssl genrsa -out private.rsa 1024
	publicKey  *rsa.PublicKey  // openssl rsa -in private.rsa -pubout > public.rsa.pub
)

func init() {
	privateBytes, err := ioutil.ReadFile("/home/terry/llavesRSA/private.rsa")
	if err != nil {
		log.Fatal("No se pudo leer el archivo privado", err)
	}

	publicBytes, err := ioutil.ReadFile("/home/terry/llavesRSA/public.rsa.pub")
	if err != nil {
		log.Fatal("No se pudo leer el archivo publico", err)
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		log.Fatal("No se pudo generar el parse privatekey")
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		log.Fatal("no se pudo generar el parse publickey", err)
	}
}

func Login(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var user, userDB model.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Fprintf(w, "Error al leer el usuario")
		return
	}
	db.Where("username = ?", user.Username).First(&userDB)

	if validateUser(user, userDB.Password) {
		user.Password = ""
		user.Role = userDB.Role
		token := GenerateJWT(user, 5)
		tokenRefresh := GenerateRefreshJWT(user, 15)
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
	} else {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Usuario o clave no válido")
	}
}

func validateUser(user model.User, password string) bool {
	match := CheckPasswordHash(user.Password, password)
	return match
}

func GenerateJWT(user model.User, timeExpire time.Duration) string {
	claims := model.Claim{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * timeExpire).Unix(),
			Issuer:    "Accesso Portafolio",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	result, err := token.SignedString(privateKey)
	if err != nil {
		fmt.Println(err)
		log.Fatal("no se pudo firmar el token", err)
	}

	//redisRefreshToken(user, token)
	return result

}

func GenerateRefreshJWT(user model.User, timeExpire time.Duration) string {
	uuid := uuid.NewV4().String()
	claims := model.ClaimRefreshToken{
		Username:    user.Username,
		RefreshUuid: uuid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * timeExpire).Unix(),
			Issuer:    "Refresh Token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	result, err := token.SignedString(privateKey)
	if err != nil {
		fmt.Println(err)
		log.Fatal("no se pudo firmar el token", err)
	}

	return result

}

func ValidateToken(w http.ResponseWriter, r *http.Request) {
	token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &model.Claim{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		switch err.(type) {
		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				fmt.Fprintf(w, "Su token ha expirado")
				return
			case jwt.ValidationErrorSignatureInvalid:
				fmt.Fprintf(w, "La firma del token no coincide")
				return
			default:
				fmt.Fprintf(w, "Su token no es valido")
				return
			}
		default:
			fmt.Fprintf(w, "Su token no es valido")
			return
		}
	}
	if token.Valid {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, "Bievenido al sistema")
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Bievenido al sistema")
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetPublicKey() *rsa.PublicKey {
	return publicKey
}
