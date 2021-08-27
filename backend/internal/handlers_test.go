package internal_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortnener/internal"

	_ "github.com/lib/pq" // Postgres Driver
)

func TestCreateShortHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	shortRequest := &internal.CreateShortRequest{URL: "lucastephens.com"}
	body, err := json.Marshal(shortRequest)
	if err != nil {
		t.Fatalf("failed to marshal json body: %v", err)
	}

	req, err := http.NewRequest("POST", "/short", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Instantiate a DAO
	dbstring := "user=user dbname=public password=password host=localhost sslmode=disable"
	driver := "postgres"

	db, err := sql.Open(driver, dbstring)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	dao := internal.NewShortPostgresDao(db, driver)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	createShort := internal.NewCreateShortHandler(dao)
	handler := http.HandlerFunc(createShort)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
