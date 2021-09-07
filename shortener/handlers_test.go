//go:build unit || all

package shortener_test

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortnener/shortener"
	"url-shortnener/test/mocks"

	"github.com/gavv/httpexpect/v2"
	"github.com/golang/mock/gomock"
)

func TestCreateShortHandler(t *testing.T) {
	type testCase struct {
		Name           string
		URL            string
		ExpectedScheme string
		ExpectedHost   string
		ExpectedPath   string
		ExpectedQuery  string
	}

	testCases := []testCase{
		{
			Name:           "Short with Nothing",
			URL:            "lucastephens.com",
			ExpectedScheme: "http",
			ExpectedHost:   "lucastephens.com",
			ExpectedPath:   "",
			ExpectedQuery:  "",
		},
		{
			Name:           "Short with Path",
			URL:            "lucastephens.com/resume.pdf",
			ExpectedScheme: "http",
			ExpectedHost:   "lucastephens.com",
			ExpectedPath:   "/resume.pdf",
			ExpectedQuery:  "",
		},
		{
			Name:           "Short with Query",
			URL:            "lucastephens.com?a=b&c=d",
			ExpectedScheme: "http",
			ExpectedHost:   "lucastephens.com",
			ExpectedPath:   "",
			ExpectedQuery:  "a=b&c=d",
		},
		{
			Name:           "Short with Everything",
			URL:            "lucastephens.com/resume.pdf?a=b&c=d",
			ExpectedScheme: "http",
			ExpectedHost:   "lucastephens.com",
			ExpectedPath:   "/resume.pdf",
			ExpectedQuery:  "a=b&c=d",
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			mock := gomock.NewController(t)
			dao := mocks.NewMockShortDAO(mock)
			dao.
				EXPECT().
				InsertShort(gomock.AssignableToTypeOf(shortener.Short{})).
				Return(nil).
				Times(1)

			createShort := shortener.NewCreateShortHandler(dao)
			handler := http.HandlerFunc(createShort)

			server := httptest.NewServer(handler)
			defer server.Close()
			e := httpexpect.New(t, server.URL)

			response := e.POST("/short").WithJSON(&shortener.CreateShortRequest{URL: test.URL}).WithHeader("Content-Type", "application/json").
				Expect().
				Status(http.StatusOK).JSON().Object()

			response.Keys().ContainsOnly("redirect_path", "scheme", "host", "path", "query")
			response.Value("redirect_path").NotNull().String()
			response.Value("scheme").NotNull().String().Equal(test.ExpectedScheme)
			response.Value("host").NotNull().String().Equal(test.ExpectedHost)
			response.Value("path").NotNull().String().Equal(test.ExpectedPath)
			response.Value("query").NotNull().String().Equal(test.ExpectedQuery)
		})
	}
}

func TestGetShortHandler(t *testing.T) {
	type testCase struct {
		Name           string
		Hash           string
		ExpectedShort  *shortener.Short
		ExpectedError  error
		ExpectedStatus int
	}

	testCases := []testCase{
		{
			Name: "Redirect",
			Hash: "c3xd4d",
			ExpectedShort: &shortener.Short{
				RedirectPath: "c3xd4d",
				Scheme:       "http",
				Host:         "github.com",
				Path:         pointerString("/soggycactus"),
				Query:        nil,
			},
			ExpectedError:  nil,
			ExpectedStatus: http.StatusMovedPermanently,
		},
		{
			Name:           "Not Found",
			Hash:           "c3xd4d",
			ExpectedShort:  nil,
			ExpectedError:  sql.ErrNoRows,
			ExpectedStatus: http.StatusNotFound,
		},
		{
			Name:           "Internal Error",
			Hash:           "c3xd4d",
			ExpectedShort:  nil,
			ExpectedError:  errors.New("internal error"),
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			mock := gomock.NewController(t)
			dao := mocks.NewMockShortDAO(mock)
			dao.
				EXPECT().
				GetShort(gomock.Any()).
				Return(test.ExpectedShort, test.ExpectedError).
				Times(1)

			getShort := shortener.NewGetShortHandler(dao)
			handler := http.HandlerFunc(getShort)

			server := httptest.NewServer(handler)
			defer server.Close()
			e := httpexpect.New(t, server.URL)

			_ = e.GET("/{short}").
				WithPath("short", test.Hash).
				WithHeader("Content-Type", "application/json").
				WithRedirectPolicy(httpexpect.DontFollowRedirects).
				Expect().
				Status(test.ExpectedStatus)
		})
	}
}

func pointerString(s string) *string {
	return &s
}
