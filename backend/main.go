package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"
	"url-shortnener/internal"

	_ "github.com/lib/pq" // Postgres Driver

	"github.com/gorilla/mux"
)

func main() {
	dbstring := os.Getenv("DBSTRING")
	driver := os.Getenv("DRIVER")

	db, err := sql.Open(driver, dbstring)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	router := mux.NewRouter()

	dao := internal.NewShortPostgresDao(db, driver)
	getShortHandler := internal.NewGetShortHandler(dao)
	createShortHandler := internal.NewCreateShortHandler(dao)

	router.HandleFunc("/{short}", getShortHandler).Methods(http.MethodGet)
	router.HandleFunc("/short", createShortHandler).Methods(http.MethodPost)

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Print("starting server on port 8080")
	log.Fatal(srv.ListenAndServe())
}
