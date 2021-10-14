package main

import (
	"encoding/json"
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
	DB.AutoMigrate(
		model.User{}, model.Company{}, model.Perfil{}, model.CatCountry{},
		model.CatServer{}, model.CatSystemOperative{}, model.CatBackEnd{},
		model.CatFrontEnd{}, model.DataBase{}, model.BackEnd{}, model.FrontEnd{},
		model.Servers{}, model.ConfConcurrency{},
	)
	return DB
}

func main() {
	app.DB = database()
	app.Router = mux.NewRouter()

	app.Router.HandleFunc("/user", app.createUser).Methods("POST")
	app.Router.HandleFunc("/user", logging(app.getAllUser)).Methods("GET")
	app.Router.HandleFunc("/user/{name}", logging(app.getUser)).Methods("GET")
	app.Router.HandleFunc("/user/{name}", app.updateUser).Methods("PUT")
	app.Router.HandleFunc("/user/{name}", app.deleteUser).Methods("DELETE")
	app.Router.HandleFunc("/login", app.login).Methods("POST")
	app.Router.HandleFunc("/logout", app.logout).Methods("POST")

	app.Router.HandleFunc("/tokenRefresh", app.tokenRefresh).Methods("POST")
	app.Router.HandleFunc("/deleteTokenRefreshRedis", app.deleteTokenRefreshRedis).Methods("POST")
	app.Router.HandleFunc("/validate", app.validate).Methods("GET")
	app.Router.HandleFunc("/validate", app.validate).Methods("POST")

	//concurrency
	app.Router.HandleFunc("/concurrency", app.concurrency).Methods("POST")
	app.Router.HandleFunc("/concurrency", app.getConcurrency).Methods("GET")
	app.Router.HandleFunc("/email", app.email).Methods("POST")

	// company
	app.Router.HandleFunc("/catalogueCountry", logging(app.catalogueCountry)).Methods("GET")
	app.Router.HandleFunc("/catalogueCompany", logging(app.catalogueCompany)).Methods("GET")
	app.Router.HandleFunc("/company", logging(app.createCompany)).Methods("POST")
	app.Router.HandleFunc("/company", logging(app.findCompany)).Methods("GET")
	app.Router.HandleFunc("/company/{companyID}", logging(app.updateCompany)).Methods("PUT")
	app.Router.HandleFunc("/company/{companyID}", logging(app.deleteCompany)).Methods("DELETE")
	app.Router.HandleFunc("/companyUpdate/{companyID}", logging(app.getCompanyUpdate)).Methods("GET")

	// perfil
	app.Router.HandleFunc("/perfil", logging(app.createPerfil)).Methods("POST")
	app.Router.HandleFunc("/perfil", logging(app.findMountedPerfil)).Methods("GET")
	app.Router.HandleFunc("/perfilFind", logging(app.findPerfil)).Methods("GET")
	app.Router.HandleFunc("/perfil/{perfilID}", logging(app.updatePerfil)).Methods("PUT")
	app.Router.HandleFunc("/perfil/{perfilID}", logging(app.deletePerfil)).Methods("DELETE")
	app.Router.HandleFunc("/mountPerfil/{perfilID}", logging(app.mountPerfil)).Methods("GET")

	http.Handle("/", app.Router)
	// db.Conexion(app.Router)

	header, methods, origin, creds := cors()
	log.Fatal(http.ListenAndServe(":3000", handlers.CORS(header, methods, origin, creds)(app.Router)))

}

func (a *App) concurrency(w http.ResponseWriter, r *http.Request) {
	handler.Concurrency(a.DB, w, r)
}

func (a *App) getConcurrency(w http.ResponseWriter, r *http.Request) {
	handler.GetConcurrency(a.DB, w, r)
}

func (a *App) email(w http.ResponseWriter, r *http.Request) {
	handler.Email(a.DB, w, r)
}

func (a *App) getCompanyUpdate(w http.ResponseWriter, r *http.Request) {
	handler.GetCompanyUpdate(a.DB, w, r)
}

func (a *App) updateCompany(w http.ResponseWriter, r *http.Request) {
	handler.UpdateCompany(a.DB, w, r)
}

func (a *App) deleteCompany(w http.ResponseWriter, r *http.Request) {
	handler.DeleteCompany(a.DB, w, r)
}

func (a *App) createPerfil(w http.ResponseWriter, r *http.Request) {
	handler.CreatePerfil(a.DB, w, r)
}

func (a *App) updatePerfil(w http.ResponseWriter, r *http.Request) {
	handler.UpdatePerfil(a.DB, w, r)
}

func (a *App) deletePerfil(w http.ResponseWriter, r *http.Request) {
	handler.DeletePerfil(a.DB, w, r)
}

func (a *App) findMountedPerfil(w http.ResponseWriter, r *http.Request) {
	handler.FindMountedPerfil(a.DB, w, r)
}

func (a *App) findPerfil(w http.ResponseWriter, r *http.Request) {
	handler.FindPerfil(a.DB, w, r)
}

func (a *App) catalogueCountry(w http.ResponseWriter, r *http.Request) {
	handler.CatalogueCountry(a.DB, w, r)
}

func (a *App) catalogueCompany(w http.ResponseWriter, r *http.Request) {
	handler.CatalogueCompany(a.DB, w, r)
}

func (a *App) mountPerfil(w http.ResponseWriter, r *http.Request) {
	handler.MountPerfil(a.DB, w, r)
}

func (a *App) createCompany(w http.ResponseWriter, r *http.Request) {
	handler.CreateCompany(a.DB, w, r)
}

func (a *App) findCompany(w http.ResponseWriter, r *http.Request) {
	handler.FindCompany(a.DB, w, r)
}

func (a *App) deleteTokenRefreshRedis(w http.ResponseWriter, r *http.Request) {
	handler.DeleteTokenRefreshRedis(a.DB, w, r)
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

func cors() (handlers.CORSOption, handlers.CORSOption, handlers.CORSOption, handlers.CORSOption) {
	enviroment := env.Env()
	header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Access-Control-Allow-Credentials", "Access-Control-Allow-Origin"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origin := handlers.AllowedOrigins([]string{enviroment["ALLOWEDORIGINS"]})
	creds := handlers.AllowCredentials()
	return header, methods, origin, creds
}

func logging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var datos Data
		datos.Validate, datos.Mensaje = middleware.ValidateToken(w, r)
		if datos.Validate {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(datos)
		}
	}
}

// func setupResponse(w *http.ResponseWriter, req *http.Request) {
// 	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
// 	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 	(*w).Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Credentials, Access-Control-Allow-Origin, withCredentials, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Origin, X-Requested-With")
// 	(*w).Header().Set("Content-Type", "text/html; charset=utf-8; charset=ascii")
// 	(*w).Header().Set("Access-Control-Allow-Credentials", "true")

// }
