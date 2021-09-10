package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"
	"url-shortnener/shortener"

	"github.com/gorilla/mux"
	"github.com/pressly/goose/v3"
	"github.com/rs/cors"
)

func main() {
	dbstring := os.Getenv("DBSTRING")
	driver := os.Getenv("DRIVER")

	db, err := sql.Open(driver, dbstring)
	if err != nil {
		log.Fatalf("failed to open to database: %v", err)
	}

	err = goose.Up(db, "migrations")
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	router := mux.NewRouter()

	dao := shortener.NewShortPostgresDao(db, driver)
	getShortHandler := shortener.NewGetShortHandler(dao)
	createShortHandler := shortener.NewCreateShortHandler(dao)

	router.HandleFunc("/{short}", getShortHandler).Methods(http.MethodGet)
	router.HandleFunc("/short", createShortHandler).Methods(http.MethodPost)
	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) { rw.WriteHeader(200) })

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5000", "https://shortener.dev"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "HEAD", "POST"},
	})

	router.Use(mux.CORSMethodMiddleware(router))

	srv := &http.Server{
		Handler:      c.Handler(router),
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Print("starting server on port 8080")
	log.Fatal(srv.ListenAndServe())
}
