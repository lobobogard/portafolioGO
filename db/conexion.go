package db

import (
	"log"
	"net/http"
	"time"
)

func Conexion(router http.Handler) {
	srv := &http.Server{
		Handler:      router,
		Addr:         "localhost:3000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
