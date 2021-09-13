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
		Name             string
		URL              string
		ExpectedScheme   string
		ExpectedHost     string
		ExpectedPath     string
		ExpectedQuery    string
		ExpectedFragment string
	}

	testCases := []testCase{
		{
			Name:             "Short with Nothing",
			URL:              "lucastephens.com",
			ExpectedScheme:   "http",
			ExpectedHost:     "lucastephens.com",
			ExpectedPath:     "",
			ExpectedQuery:    "",
			ExpectedFragment: "",
		},
		{
			Name:             "Short with Path",
			URL:              "lucastephens.com/resume.pdf",
			ExpectedScheme:   "http",
			ExpectedHost:     "lucastephens.com",
			ExpectedPath:     "/resume.pdf",
			ExpectedQuery:    "",
			ExpectedFragment: "",
		},
		{
			Name:             "Short with Query",
			URL:              "lucastephens.com?a=b&c=d",
			ExpectedScheme:   "http",
			ExpectedHost:     "lucastephens.com",
			ExpectedPath:     "",
			ExpectedQuery:    "a=b&c=d",
			ExpectedFragment: "",
		},
		{
			Name:             "Short with Fragment",
			URL:              "lucastephens.com#info",
			ExpectedScheme:   "http",
			ExpectedHost:     "lucastephens.com",
			ExpectedPath:     "",
			ExpectedQuery:    "",
			ExpectedFragment: "info",
		},
		{
			Name:             "Short with Everything",
			URL:              "lucastephens.com/resume.pdf?a=b&c=d#info",
			ExpectedScheme:   "http",
			ExpectedHost:     "lucastephens.com",
			ExpectedPath:     "/resume.pdf",
			ExpectedQuery:    "a=b&c=d",
			ExpectedFragment: "info",
		},
		{
			Name:             "Short with Host Path Fragment",
			URL:              "https://mail.google.com/mail/u/2/#inbox",
			ExpectedScheme:   "https",
			ExpectedHost:     "mail.google.com",
			ExpectedPath:     "/mail/u/2/",
			ExpectedQuery:    "",
			ExpectedFragment: "inbox",
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			mock := gomock.NewController(t)
			dao := mocks.NewMockShortDAO(mock)
			dao.
				EXPECT().
				InsertShort(gomock.Any(), gomock.AssignableToTypeOf(shortener.Short{})).
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

			response.Keys().ContainsOnly("redirect_path", "scheme", "host", "path", "query", "fragment")
			response.Value("redirect_path").NotNull().String()
			response.Value("scheme").NotNull().String().Equal(test.ExpectedScheme)
			response.Value("host").NotNull().String().Equal(test.ExpectedHost)
			response.Value("path").NotNull().String().Equal(test.ExpectedPath)
			response.Value("query").NotNull().String().Equal(test.ExpectedQuery)
			response.Value("fragment").NotNull().String().Equal(test.ExpectedFragment)
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
				GetShort(gomock.Any(), gomock.Any()).
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
