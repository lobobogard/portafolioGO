package authentication

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/request"
	"github.com/joho/godotenv"
	"github.com/portafolioLP/env"
	"github.com/portafolioLP/model"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//create keys privadas y publicas minuto 26
var (
	privateKey *rsa.PrivateKey // openssl genrsa -out private.rsa 2048 / 1024
	publicKey  *rsa.PublicKey  // openssl rsa -in private.rsa -pubout > public.rsa.pub
)

type TokenRefresh struct {
	Token string
}

func Env() map[string]string {
	var Env = make(map[string]string)
	err := godotenv.Load(filepath.Join("/root/", ".env"))

	if err != nil {
		log.Fatal("Error loading .env authentificaction.go ", err)
	}

	envPath, _ := filepath.Abs("/root/.env")
	Env, err = godotenv.Read(envPath)
	if err != nil {
		log.Fatal("error reading .env authentificaction.go ", err)
	}

	return Env

}

func init() {
	Env := Env()                                            // error testing
	privateBytes, err := ioutil.ReadFile(Env["PRIVATEKEY"]) // error testing
	// privateBytes, err := ioutil.ReadFile("/home/terry/llavesRSA/private.rsa") // for testing
	// privateBytes, err := ioutil.ReadFile("/usr/llavesRSA/private.rsa") // for testing
	if err != nil {
		log.Fatal("Could not read private file", err)
	}

	publicBytes, err := ioutil.ReadFile(Env["PUBLICKEY"]) // error testing
	// publicBytes, err := ioutil.ReadFile("/home/terry/llavesRSA/public.rsa.pub") // for testing
	// publicBytes, err := ioutil.ReadFile("/usr/llavesRSA/public.rsa.pub") // for testing
	if err != nil {
		log.Fatal("Could not read the public file", err)
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		log.Fatal("Could not generate privatekey parse")
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		log.Fatal("could not generate the publickey parse", err)
	}
}

func Login(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var user, userDB model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Fprintf(w, "Error reading user")
		return
	}
	db.Where("username = ?", user.Username).First(&userDB)

	if ValidateUserLogin(user, userDB.Password) {
		user.Password = ""
		user.Role = "USER" // userDB.Role
		token, err := GenerateJWT(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(err.Error()))
			return
		}

		_, err = GenerateRefreshJWT(w, user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(err.Error()))
			return
		}

		result := model.ResponseToken{Token: token}
		jsonResult, err := json.Marshal(result)
		if err != nil {
			jsonResult, _ := json.Marshal("Internal error in the system")
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResult)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResult)
	} else {
		jsonResult, _ := json.Marshal("Invalid username or password")
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResult)
	}
}

func GenerateJWT(user model.User) (string, error) {
	claims := model.Claim{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 1).Unix(), // 1 minute
			Issuer:    "access token",
			Id:        user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	result, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return result, err

}

func GenerateRefreshJWT(w http.ResponseWriter, user model.User) (string, error) {
	uuid := uuid.NewV4().String()

	claims := model.ClaimRefreshToken{
		Username:    user.Username,
		RefreshUuid: uuid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 20).Unix(), // 20 minute
			Issuer:    "Refresh Token",
		},
	}

	tokenRefresh := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	result, err := tokenRefresh.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	// generate tokenRefresh cookie
	Env := env.Env()

	http.SetCookie(w, &http.Cookie{
		Name:     "tokenRefresh",
		Path:     "/",
		Domain:   Env["DOMAINCOOKIE"],
		Value:    result,
		Expires:  time.Now().Add(8 * time.Hour),
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
		Secure:   true,
	})

	// generate tokenRefresh in redis
	model.RedisRefreshToken(user, result)

	return result, err

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
				fmt.Fprintf(w, "Your token has expired")
				return
			case jwt.ValidationErrorSignatureInvalid:
				fmt.Fprintf(w, "Token signature does not match")
				return
			default:
				fmt.Fprintf(w, "Your token is not valid")
				return
			}
		default:
			fmt.Fprintf(w, "Your token is not valid")
			return
		}
	}
	if token.Valid {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, "You welcome the system")
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorized")
	}
}

func ValidateUserLogin(user model.User, password string) bool {
	if user.Username == "" {
		return false
	}
	match := CheckPasswordHash(user.Password, password)
	return match
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetPublicKey() *rsa.PublicKey {
	return publicKey
}

func Logout(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("tokenRefresh")
	if err != nil {
		fmt.Printf("Not cookie for delete status ok\n")
		return
	}
	DeleteTokenRefreshRedisCookie(cookie.Value)

	c := &http.Cookie{
		Name:     "tokenRefresh",
		Value:    "",
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
	}

	http.SetCookie(w, c)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResult, _ := json.Marshal("Logout Success")
	w.Write(jsonResult)
}

func DeleteTokenRefreshRedisCookie(tokenStr string) {
	token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return GetPublicKey(), nil
	})
	claims, _ := token.Claims.(jwt.MapClaims)

	username := fmt.Sprintf("%v", claims["jti"])
	model.RedisDeleteRefreshToken(username)

}
