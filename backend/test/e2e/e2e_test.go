//go:build e2e || all

package e2e

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortnener/shortener"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // Postgres Driver

	"github.com/gavv/httpexpect/v2"
)

func TestCreateAndGet(t *testing.T) {
	shortRequest := &shortener.CreateShortRequest{URL: "lucastephens.com"}

	dbstring := "user=user dbname=public password=password host=localhost sslmode=disable"
	driver := "postgres"

	db, err := sql.Open(driver, dbstring)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	dao := shortener.NewShortPostgresDao(db, driver)
	getShortHandler := shortener.NewGetShortHandler(dao)
	createShortHandler := shortener.NewCreateShortHandler(dao)

	router := mux.NewRouter()
	router.HandleFunc("/{short}", getShortHandler).Methods(http.MethodGet)
	router.HandleFunc("/short", createShortHandler).Methods(http.MethodPost)

	server := httptest.NewServer(router)
	defer server.Close()
	e := httpexpect.New(t, server.URL)
	postResponse := e.POST("/short").WithJSON(shortRequest).WithHeader("Content-Type", "application/json").
		Expect().
		Status(http.StatusOK).JSON().Object()

	postResponse.Keys().ContainsOnly("redirect_path", "scheme", "host", "path", "query")
	redirect_path := postResponse.Value("redirect_path").NotNull().String()
	postResponse.Value("scheme").NotNull().String().Equal("http")
	postResponse.Value("host").NotNull().String().Equal("lucastephens.com")
	postResponse.Value("path").NotNull().String().Equal("")
	postResponse.Value("query").NotNull().String().Equal("")

	getResponse := e.GET("/{short}").
		WithPath("short", redirect_path.Raw()).
		WithHeader("Content-Type", "application/json").
		WithRedirectPolicy(httpexpect.DontFollowRedirects).
		Expect().
		Status(http.StatusMovedPermanently)

	getResponse.Body().NotEmpty().Contains("http://lucastephens.com")
}
