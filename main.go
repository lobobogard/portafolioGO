package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/portafolioLP/authentication"
	"github.com/portafolioLP/db"
	"github.com/portafolioLP/env"
	"github.com/portafolioLP/handler"
	"github.com/portafolioLP/libs"
	"github.com/portafolioLP/middleware"
	"github.com/portafolioLP/model"
	"gorm.io/gorm"
	// "github.com/portafolioLP/login"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

var app App

func database() {
	Env := env.Env()
	dbConfig := libs.Configure(Env)
	app.DB = dbConfig.InitMysqlDB()
	app.DB.AutoMigrate(model.User{}, model.Company{}, model.Perfil{})
}

// func generateJWT() {
// 	token := jwt.New(jwt.SigningMethodHS256)
// 	fmt.Println("token", token)
// }

func main() {
	// generateJWT()
	database()
	app.Router = mux.NewRouter()

	app.Router.HandleFunc("/user", app.createUser).Methods("POST")
	app.Router.HandleFunc("/user", logging(app.getAllUser)).Methods("GET")
	app.Router.HandleFunc("/user/{name}", app.getUser).Methods("GET")
	app.Router.HandleFunc("/user/{name}", app.updateUser).Methods("PUT")
	app.Router.HandleFunc("/user/{name}", app.deleteUser).Methods("DELETE")
	app.Router.HandleFunc("/login", app.login).Methods("POST")
	app.Router.HandleFunc("/tokenRefresh", app.tokenRefresh).Methods("POST")
	app.Router.HandleFunc("/validate", app.validate).Methods("GET")
	http.Handle("/", app.Router)
	db.Conexion(app.Router)

}

func logging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := middleware.ValidateToken(w, r)
		if resp {
			fmt.Println(resp)
			next.ServeHTTP(w, r)
		}
	}
}

func (a *App) tokenRefresh(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	handler.ValidateTokenRefresh(a.DB, w, r)
}

func (a *App) validate(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	authentication.ValidateToken(w, r)
}

func (a *App) login(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	authentication.Login(a.DB, w, r)
}

func (a *App) getAllUser(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	handler.GetAllUsers(a.DB, w, r)
}

func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	handler.User(a.DB, w, r)
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	handler.GetUser(a.DB, w, r)
}

func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	handler.UpdateUser(a.DB, w, r)
}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	handler.DeleteUser(a.DB, w, r)
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
