package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"l24.dev/shortener"

	"github.com/gorilla/mux"
	"github.com/pressly/goose/v3"
	"github.com/rs/cors"
)

func getDBString() string {
	// These environment variables are set automatically by Cloud Run
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, name)
}

func main() {
	var dbstring string

	// DBSTRING is filled out in local dev
	dbstring = os.Getenv("DBSTRING")
	// If not present, user the env vars set by Cloud Run
	if dbstring == "" {
		log.Print("DBSTRING not present, looking for Cloud Run env vars")
		dbstring = getDBString()
	}
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
