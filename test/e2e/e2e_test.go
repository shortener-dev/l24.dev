//go:build e2e || all

package e2e

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"l24.dev/shortener"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // Postgres Driver

	"github.com/gavv/httpexpect/v2"
)

func TestCreateAndGet(t *testing.T) {
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

	type testCase struct {
		Name             string
		Request          shortener.CreateShortRequest
		ExpectedScheme   string
		ExpectedHost     string
		ExpectedPath     string
		ExpectedQuery    string
		ExpectedFragment string
		ExpectedGetBody  string
	}

	testCases := []testCase{
		{
			Name:             "Host Only",
			Request:          shortener.CreateShortRequest{URL: "lucastephens.com"},
			ExpectedScheme:   "http",
			ExpectedHost:     "lucastephens.com",
			ExpectedPath:     "",
			ExpectedQuery:    "",
			ExpectedFragment: "",
			ExpectedGetBody:  "http://lucastephens.com",
		},
		{
			Name:             "Host Path Fragment",
			Request:          shortener.CreateShortRequest{URL: "https://mail.google.com/mail/u/2/#inbox"},
			ExpectedScheme:   "https",
			ExpectedHost:     "mail.google.com",
			ExpectedPath:     "/mail/u/2/",
			ExpectedQuery:    "",
			ExpectedFragment: "inbox",
			ExpectedGetBody:  "https://mail.google.com/mail/u/2/#inbox",
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			postResponse := e.POST("/short").WithJSON(test.Request).WithHeader("Content-Type", "application/json").
				Expect().
				Status(http.StatusOK).JSON().Object()

			postResponse.Keys().ContainsOnly("redirect_path", "scheme", "host", "path", "query", "fragment")
			redirect_path := postResponse.Value("redirect_path").NotNull().String()
			postResponse.Value("scheme").NotNull().String().Equal(test.ExpectedScheme)
			postResponse.Value("host").NotNull().String().Equal(test.ExpectedHost)
			postResponse.Value("path").NotNull().String().Equal(test.ExpectedPath)
			postResponse.Value("query").NotNull().String().Equal(test.ExpectedQuery)
			postResponse.Value("fragment").NotNull().String().Equal(test.ExpectedFragment)

			getResponse := e.GET("/{short}").
				WithPath("short", redirect_path.Raw()).
				WithHeader("Content-Type", "application/json").
				WithRedirectPolicy(httpexpect.DontFollowRedirects).
				Expect().
				Status(http.StatusMovedPermanently)

			getResponse.Body().NotEmpty().Contains(test.ExpectedGetBody)
		})
	}

}
