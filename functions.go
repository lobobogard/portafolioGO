package main

import (
	"net/http"

	"github.com/portafolioLP/authentication"
	"github.com/portafolioLP/handler"
)

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

func (a *App) login(w http.ResponseWriter, r *http.Request) {
	authentication.Login(a.DB, w, r)
}

func (a *App) logout(w http.ResponseWriter, r *http.Request) {
	authentication.Logout(a.DB, w, r)
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

func (a *App) createCompany(w http.ResponseWriter, r *http.Request) {
	handler.CreateCompany(a.DB, w, r)
}

func (a *App) findCompany(w http.ResponseWriter, r *http.Request) {
	handler.FindCompany(a.DB, w, r)
}

func (a *App) catalogueCountry(w http.ResponseWriter, r *http.Request) {
	handler.CatalogueCountry(a.DB, w, r)
}

func (a *App) catalogueCompany(w http.ResponseWriter, r *http.Request) {
	handler.CatalogueCompany(a.DB, w, r)
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

func (a *App) mountPerfil(w http.ResponseWriter, r *http.Request) {
	handler.MountPerfil(a.DB, w, r)
}

func (a *App) estadistic(w http.ResponseWriter, r *http.Request) {
	handler.Estadistic(a.DB, w, r)
}

func (a *App) mountEstadistic(w http.ResponseWriter, r *http.Request) {
	handler.MountEstadistic(a.DB, w, r)
}

func (a *App) tokenRefresh(w http.ResponseWriter, r *http.Request) {
	handler.ValidateTokenRefresh(a.DB, w, r)
}

func (a *App) deleteTokenRefreshRedis(w http.ResponseWriter, r *http.Request) {
	handler.DeleteTokenRefreshRedis(a.DB, w, r)
}

func (a *App) validate(w http.ResponseWriter, r *http.Request) {
	authentication.ValidateToken(w, r)
}

func (a *App) email(w http.ResponseWriter, r *http.Request) {
	handler.Email(a.DB, w, r)
}

func (a *App) concurrency(w http.ResponseWriter, r *http.Request) {
	handler.Concurrency(a.DB, w, r)
}

func (a *App) getConcurrency(w http.ResponseWriter, r *http.Request) {
	handler.GetConcurrency(a.DB, w, r)
}
