package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/portafolioLP/corsConfig"
	"github.com/portafolioLP/env"
	"github.com/portafolioLP/libs"
	"github.com/portafolioLP/middleware"
	"github.com/portafolioLP/model"
	"gorm.io/gorm"
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

	// example rest api
	app.Router.HandleFunc("/user", app.createUser).Methods("POST")
	app.Router.HandleFunc("/user", middleware.Logging(app.getAllUser)).Methods("GET")
	app.Router.HandleFunc("/user/{name}", middleware.Logging(app.getUser)).Methods("GET")
	app.Router.HandleFunc("/user/{name}", app.updateUser).Methods("PUT")
	app.Router.HandleFunc("/user/{name}", app.deleteUser).Methods("DELETE")

	// login
	app.Router.HandleFunc("/login", app.login).Methods("POST")
	app.Router.HandleFunc("/logout", app.logout).Methods("POST")

	// company
	app.Router.HandleFunc("/catalogueCountry", middleware.Logging(app.catalogueCountry)).Methods("GET")
	app.Router.HandleFunc("/catalogueCompany", middleware.Logging(app.catalogueCompany)).Methods("GET")
	app.Router.HandleFunc("/company", middleware.Logging(app.createCompany)).Methods("POST")
	app.Router.HandleFunc("/company", middleware.Logging(app.findCompany)).Methods("GET")
	app.Router.HandleFunc("/company/{companyID}", middleware.Logging(app.updateCompany)).Methods("PUT")
	app.Router.HandleFunc("/company/{companyID}", middleware.Logging(app.deleteCompany)).Methods("DELETE")
	app.Router.HandleFunc("/companyUpdate/{companyID}", middleware.Logging(app.getCompanyUpdate)).Methods("GET")

	// perfil
	app.Router.HandleFunc("/perfil", middleware.Logging(app.createPerfil)).Methods("POST")
	app.Router.HandleFunc("/perfil", middleware.Logging(app.findMountedPerfil)).Methods("GET")
	app.Router.HandleFunc("/perfilFind", middleware.Logging(app.findPerfil)).Methods("GET")
	app.Router.HandleFunc("/perfil/{perfilID}", middleware.Logging(app.updatePerfil)).Methods("PUT")
	app.Router.HandleFunc("/perfil/{perfilID}", middleware.Logging(app.deletePerfil)).Methods("DELETE")
	app.Router.HandleFunc("/mountPerfil/{perfilID}", middleware.Logging(app.mountPerfil)).Methods("GET")

	// estadistic
	app.Router.HandleFunc("/estadistic", middleware.Logging(app.estadistic)).Methods("GET")
	app.Router.HandleFunc("/mountEstadistic", middleware.Logging(app.mountEstadistic)).Methods("GET")

	// token
	app.Router.HandleFunc("/tokenRefresh", app.tokenRefresh).Methods("POST")
	app.Router.HandleFunc("/deleteTokenRefreshRedis", app.deleteTokenRefreshRedis).Methods("POST")

	// validation
	app.Router.HandleFunc("/validate", app.validate).Methods("GET")
	app.Router.HandleFunc("/validate", app.validate).Methods("POST")

	//concurrency
	app.Router.HandleFunc("/concurrency", app.concurrency).Methods("POST")
	app.Router.HandleFunc("/concurrency", app.getConcurrency).Methods("GET")
	app.Router.HandleFunc("/email", app.email).Methods("POST")

	http.Handle("/", app.Router)

	header, methods, origin, creds := corsConfig.Cors()
	Env := env.Env()
	log.Fatal(http.ListenAndServe(Env["PORTS"], handlers.CORS(header, methods, origin, creds)(app.Router)))

}
