//go:build unit || all

package shortener_test

import (
	"context"
	"database/sql/driver"
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
		ShouldFail    bool
	}

	testCases := []testCase{
		{
			Name:          "Short with Path Only",
			ExpectedQuery: shortener.InsertShortWithPathQuery,
			Scheme:        "http",
			Host:          "github.com",
			Path:          "/DATA-DOG/go-sqlmock",
			Query:         "",
			ShouldFail:    false,
		},
		{
			Name:          "Short with Query Only",
			ExpectedQuery: shortener.InsertShortWithQueryQuery,
			Scheme:        "http",
			Host:          "github.com",
			Path:          "",
			Query:         "test=value",
			ShouldFail:    false,
		},
		{
			Name:          "Short with Everything",
			ExpectedQuery: shortener.InsertShortWithAllQuery,
			Scheme:        "http",
			Host:          "github.com",
			Path:          "/soggycactus",
			Query:         "test=value",
			ShouldFail:    false,
		},
		{
			Name:          "Short with Nothing",
			ExpectedQuery: shortener.InsertShortQuery,
			Scheme:        "http",
			Host:          "github.com",
			Path:          "",
			Query:         "",
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
			}

			var args []driver.Value
			args = append(args, "test")
			args = append(args, test.Scheme)
			args = append(args, test.Host)

			if test.Path != "" {
				args = append(args, test.Path)
			}

			if test.Query != "" {
				args = append(args, test.Query)
			}

			mock.ExpectBegin()
			mock.ExpectExec(regexp.QuoteMeta(test.ExpectedQuery)).
				WithArgs(args...).
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
