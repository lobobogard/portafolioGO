package corsConfig

import (
	"github.com/gorilla/handlers"
	"github.com/portafolioLP/env"
)

func Cors() (handlers.CORSOption, handlers.CORSOption, handlers.CORSOption, handlers.CORSOption) {
	enviroment := env.Env()
	header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Access-Control-Allow-Credentials", "Access-Control-Allow-Origin"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origin := handlers.AllowedOrigins([]string{enviroment["ALLOWEDORIGINS"]})
	creds := handlers.AllowCredentials()
	return header, methods, origin, creds
}
