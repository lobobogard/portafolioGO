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

type TokenRefresh struct {
	Token string
}

func init() {
	privateBytes, err := ioutil.ReadFile("/home/terry/llavesRSA/private.rsa")
	if err != nil {
		log.Fatal("Could not read private file", err)
	}

	publicBytes, err := ioutil.ReadFile("/home/terry/llavesRSA/public.rsa.pub")
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

	if validateUser(user, userDB.Password) {
		user.Password = ""
		user.Role = userDB.Role
		token := GenerateJWT(user)
		GenerateRefreshJWT(w, user)

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

func validateUser(user model.User, password string) bool {
	if user.Username == "" {
		return false
	}

	match := CheckPasswordHash(user.Password, password)
	return match
}

func GenerateJWT(user model.User) string {
	claims := model.Claim{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
			Issuer:    "Access Token",
			Id:        user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	result, err := token.SignedString(privateKey)
	if err != nil {
		fmt.Println(err)
		log.Fatal("the token could not be signed", err)
	}

	return result

}

func GenerateRefreshJWT(w http.ResponseWriter, user model.User) string {
	uuid := uuid.NewV4().String()

	claims := model.ClaimRefreshToken{
		Username:    user.Username,
		RefreshUuid: uuid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			Issuer:    "Refresh Token",
		},
	}

	tokenRefresh := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	result, err := tokenRefresh.SignedString(privateKey)
	if err != nil {
		fmt.Println(err)
		log.Fatal("the token could not be signed", err)
	}

	// generate tokenRefresh cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "tokenRefresh",
		Value:    result,
		Expires:  time.Now().Add(8 * time.Hour),
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
	})

	// generate tokenRefresh in redis
	model.RedisRefreshToken(user, result)

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

func Logout(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("tokenRefresh")
	if err != nil {
		fmt.Printf("Cant find cookie :/\r\n")
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
	fmt.Println(username)
	model.RedisDeleteRefreshToken(username)

}

// func DeleteTokenRefreshRedis(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		log.Printf("Reading body fail")
// 	}

// 	tokenBody := TokenRefresh{}
// 	err = json.Unmarshal(body, &tokenBody)
// 	if err != nil {
// 		log.Printf("Reading body failed: %s", err)
// 		return
// 	}

// 	tokenStr := tokenBody.Token
// 	token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
// 		return GetPublicKey(), nil
// 	})
// 	claims, _ := token.Claims.(jwt.MapClaims)

// 	username := fmt.Sprintf("%v", claims["jti"])
// 	// fmt.Println(username)
// 	model.RedisDeleteRefreshToken(username)

// }
