package authentication_test

import (
	"net/http/httptest"
	"testing"

	"github.com/portafolioLP/authentication"
	"github.com/portafolioLP/model"
)

func TestValidateUserLogin(t *testing.T) {

	tt := []struct {
		name       string
		username   string
		password   string
		passwordDB string
		result     bool
	}{
		{name: "user correct", username: "username", password: "123456", passwordDB: "$2a$14$Owa2UrjDQzOQ3GIR9uWA7usrslZyYvZndXVA.iPlLgvvgHuSCyf26", result: true},
		{name: "user not username", username: "", password: "123456", passwordDB: "$2a$14$Owa2UrjDQzOQ3GIR9uWA7usrslZyYvZndXVA.iPlLgvvgHuSCyf26", result: false},
		{name: "user not password", username: "username", password: "", passwordDB: "$2a$14$Owa2UrjDQzOQ3GIR9uWA7usrslZyYvZndXVA.iPlLgvvgHuSCyf26", result: false},
		{name: "user incorrect passwordDB", username: "username", password: "123456", passwordDB: "$2a$14$ktFLES5rADcJ4bDzv0/RHuAMm0oaX8yaAuR9NpG1WfoY6WeJtsn", result: false},
		{name: "user not username and password", username: "", password: "", passwordDB: "$2a$14$Owa2UrjDQzOQ3GIR9uWA7usrslZyYvZndXVA.iPlLgvvgHuSCyf26", result: false},
		{name: "user no data", username: "", password: "", passwordDB: "", result: false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			user := model.User{Username: tc.username, Password: tc.password}
			userDB := model.User{Password: tc.passwordDB}
			result := authentication.ValidateUserLogin(user, userDB.Password)
			if result != tc.result {
				t.Errorf("ValidateUserLogin is incorrect values %v", result)
			}
		})
	}

}

func TestGenerateJWT(t *testing.T) {
	user := model.User{Username: "terry", Password: "lobo", Role: "admin"}
	_, err := authentication.GenerateJWT(user)
	if err != nil {
		t.Errorf("GenerateJWT is incorrect error %v", err)
	}
}

func TestGenerateRefreshJWT(t *testing.T) {
	rec := httptest.NewRecorder()
	user := model.User{Username: "terry", Password: "lobo", Role: "admin"}

	_, err := authentication.GenerateRefreshJWT(rec, user)
	if err != nil {
		t.Errorf("GenerateJWT is incorrect error %v", err)
	}

}
