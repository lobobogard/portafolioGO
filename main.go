package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func Prueba(w http.ResponseWriter, r *http.Request) {
	// login.Create()
	// login.Delete()
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func main() {
	addEnvs()
	router := mux.NewRouter()
	router.HandleFunc("/", Prueba).Methods("GET")
	// r.HandleFunc("/articles", ArticlesHandler)
	http.Handle("/", router)

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
