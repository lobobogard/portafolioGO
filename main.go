package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/portafolioLP/authentication"
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

type Data struct {
	Mensaje  string
	Validate bool
}

// gorm and mux
var app App

func database() *gorm.DB {
	Env := env.Env()
	dbConfig := libs.Configure(Env)
	DB := dbConfig.InitMysqlDB()
	DB.AutoMigrate(model.User{}, model.Company{}, model.Perfil{})
	return DB
}

func main() {
	app.DB = database()
	app.Router = mux.NewRouter()

	app.Router.HandleFunc("/user", app.createUser).Methods("POST")
	app.Router.HandleFunc("/user", logging(app.getAllUser)).Methods("GET")
	app.Router.HandleFunc("/user/{name}", app.getUser).Methods("GET")
	app.Router.HandleFunc("/user/{name}", app.updateUser).Methods("PUT")
	app.Router.HandleFunc("/user/{name}", app.deleteUser).Methods("DELETE")
	app.Router.HandleFunc("/login", app.login).Methods("POST")
	app.Router.HandleFunc("/logout", app.logout).Methods("POST")
	app.Router.HandleFunc("/validate", app.validate).Methods("GET")

	app.Router.HandleFunc("/tokenRefresh", app.tokenRefresh).Methods("POST")
	http.Handle("/", app.Router)
	// db.Conexion(app.Router)

	header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Access-Control-Allow-Credentials", "Access-Control-Allow-Origin"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origin := handlers.AllowedOrigins([]string{"http://localhost:8080"})
	creds := handlers.AllowCredentials()
	log.Fatal(http.ListenAndServe(":3000", handlers.CORS(header, methods, origin, creds)(app.Router)))

}

func logging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var datos Data
		datos.Validate, datos.Mensaje = middleware.ValidateToken(w, r)
		if datos.Validate {
			fmt.Println(datos.Mensaje)
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(datos)
		}
	}
}

func (a *App) tokenRefresh(w http.ResponseWriter, r *http.Request) {
	handler.ValidateTokenRefresh(a.DB, w, r)
}

func (a *App) validate(w http.ResponseWriter, r *http.Request) {
	authentication.ValidateToken(w, r)
}

func (a *App) login(w http.ResponseWriter, r *http.Request) {
	// setupResponse(&w, r)
	authentication.Login(a.DB, w, r)
}

func (a *App) logout(w http.ResponseWriter, r *http.Request) {
	authentication.Logout(a.DB, w, r)
}

func (a *App) getAllUser(w http.ResponseWriter, r *http.Request) {
	handler.GetAllUsers(a.DB, w, r)
}

func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	handler.User(a.DB, w, r)
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	handler.GetUser(a.DB, w, r)
}

func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	handler.UpdateUser(a.DB, w, r)
}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	handler.DeleteUser(a.DB, w, r)
}

// func setupResponse(w *http.ResponseWriter, req *http.Request) {
// 	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
// 	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 	(*w).Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Credentials, Access-Control-Allow-Origin, withCredentials, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Origin, X-Requested-With")
// 	(*w).Header().Set("Content-Type", "text/html; charset=utf-8; charset=ascii")
// 	(*w).Header().Set("Access-Control-Allow-Credentials", "true")

// }
