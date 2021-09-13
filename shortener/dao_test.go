//go:build unit || all

package shortener_test

import (
	"context"
	"regexp"
	"testing"
	"url-shortnener/shortener"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestInsertShort(t *testing.T) {
	type testCase struct {
		Name          string
		ExpectedQuery string
		Scheme        string
		Host          string
		Path          string
		Query         string
		Fragment      string
		ShouldFail    bool
	}

	testCases := []testCase{
		{
			Name:          "Short with Path Only",
			ExpectedQuery: "INSERT INTO urls (redirect_path, scheme, host, path) VALUES ('test', 'http', 'github.com', '/DATA-DOG/go-sqlmock')",
			Scheme:        "http",
			Host:          "github.com",
			Path:          "/DATA-DOG/go-sqlmock",
			Query:         "",
			Fragment:      "",
			ShouldFail:    false,
		},
		{
			Name:          "Short with Query Only",
			ExpectedQuery: "INSERT INTO urls (redirect_path, scheme, host, query) VALUES ('test', 'http', 'github.com', 'test=value')",
			Scheme:        "http",
			Host:          "github.com",
			Path:          "",
			Query:         "test=value",
			Fragment:      "",
			ShouldFail:    false,
		},
		{
			Name:          "Short with Fragment Only",
			ExpectedQuery: "INSERT INTO urls (redirect_path, scheme, host, fragment) VALUES ('test', 'http', 'github.com', 'info')",
			Scheme:        "http",
			Host:          "github.com",
			Path:          "",
			Query:         "",
			Fragment:      "info",
			ShouldFail:    false,
		},
		{
			Name:          "Short with Path & Fragment",
			ExpectedQuery: "INSERT INTO urls (redirect_path, scheme, host, path, fragment) VALUES ('test', 'http', 'github.com', '/soggycactus', 'info')",
			Scheme:        "http",
			Host:          "github.com",
			Path:          "/soggycactus",
			Query:         "",
			Fragment:      "info",
			ShouldFail:    false,
		},
		{
			Name:          "Short with Query & Fragment",
			ExpectedQuery: "INSERT INTO urls (redirect_path, scheme, host, query, fragment) VALUES ('test', 'http', 'github.com', 'test=value', 'info')",
			Scheme:        "http",
			Host:          "github.com",
			Path:          "",
			Query:         "test=value",
			Fragment:      "info",
			ShouldFail:    false,
		},
		{
			Name:          "Short with Everything",
			ExpectedQuery: "INSERT INTO urls (redirect_path, scheme, host, path, query, fragment) VALUES ('test', 'http', 'github.com', '/soggycactus', 'test=value', 'info')",
			Scheme:        "http",
			Host:          "github.com",
			Path:          "/soggycactus",
			Query:         "test=value",
			Fragment:      "info",
			ShouldFail:    false,
		},
		{
			Name:          "Short with Nothing",
			ExpectedQuery: "INSERT INTO urls (redirect_path, scheme, host) VALUES ('test', 'http', 'github.com')",
			Scheme:        "http",
			Host:          "github.com",
			Path:          "",
			Query:         "",
			Fragment:      "",
			ShouldFail:    false,
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			short := shortener.Short{
				RedirectPath: "test",
				Scheme:       test.Scheme,
				Host:         test.Host,
				Path:         &test.Path,
				Query:        &test.Query,
				Fragment:     &test.Fragment,
			}

			mock.ExpectBegin()
			mock.ExpectExec(regexp.QuoteMeta(test.ExpectedQuery)).
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			dao := shortener.NewShortPostgresDao(db, "postgres")
			err = dao.InsertShort(context.Background(), short)
			if err != nil {
				t.Logf("failed to insert: %v", err)
			}

			assert.Equal(t, test.ShouldFail, err != nil, "ShouldFail is %v, got %v", test.ShouldFail, err != nil)
			assert.Nil(t, mock.ExpectationsWereMet(), "mock expectations should be met")
		})
	}
}
