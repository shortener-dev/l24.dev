//go:build unit || all

package shortener_test

import (
	"strings"
	"testing"
	"url-shortnener/shortener"

	"github.com/stretchr/testify/assert"
)

func TestDuplicateURLsAreUnique(t *testing.T) {
	short1, err := shortener.NewShort("http", "lucastephens.com", "resume.pdf", "")
	if err != nil {
		t.Fatalf("failed to create short: %v", err)
	}

	short2, err := shortener.NewShort("http", "lucastephens.com", "resume.pdf", "")
	if err != nil {
		t.Fatalf("failed to create short: %v", err)
	}

	assert.NotEqual(t, "", short1.RedirectPath, "short1 should have redirect path")
	assert.NotEqual(t, "", short2.RedirectPath, "short2 should have redirect path")
	assert.Equal(t, "http", short1.Scheme, "short1 scheme should match")
	assert.Equal(t, "http", short2.Scheme, "short2 scheme should match")
	assert.Equal(t, "lucastephens.com", short1.Host, "short1 host should match")
	assert.Equal(t, "lucastephens.com", short2.Host, "short2 host should match")
	assert.Equal(t, "/resume.pdf", *short1.Path, "short1 host should match")
	assert.Equal(t, "/resume.pdf", *short2.Path, "short2 path should match")
	assert.Equal(t, "", *short1.Query, "short1 query should match")
	assert.Equal(t, "", *short2.Query, "short2 query should match")
	assert.NotEqual(t, short1.RedirectPath, short2.RedirectPath, "shorts should have different hashes")
	assert.Equal(t, short1.RawURL(), short2.RawURL(), "shorts should have same raw urls")
	assert.Equal(t, "http://lucastephens.com/resume.pdf", short1.RawURL(), "raw url should be correct")
}

func TestNewShort(t *testing.T) {
	type testCase struct {
		Name   string
		Scheme string
		Host   string
		Path   string
		Query  string
		RawURL string
	}

	testCases := []testCase{
		{
			Name:   "Short with Everything",
			Scheme: "http",
			Host:   "lucastephens.com",
			Path:   "resume.pdf",
			Query:  "?a=b&c=d",
			RawURL: "http://lucastephens.com/resume.pdf?a=b&c=d",
		},
		{
			Name:   "Short with Path Only",
			Scheme: "http",
			Host:   "lucastephens.com",
			Path:   "resume.pdf",
			Query:  "",
			RawURL: "http://lucastephens.com/resume.pdf",
		},
		{
			Name:   "Short with Query Only",
			Scheme: "http",
			Host:   "lucastephens.com",
			Path:   "",
			Query:  "?a=b&c=d",
			RawURL: "http://lucastephens.com?a=b&c=d",
		},
		{
			Name:   "Short with Nothing",
			Scheme: "http",
			Host:   "lucastephens.com",
			Path:   "",
			Query:  "",
			RawURL: "http://lucastephens.com",
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			short, err := shortener.NewShort(test.Scheme, test.Host, test.Path, test.Query)
			if err != nil {
				t.Fatalf("failed to create short: %v", err)
			}

			expectedPath := test.Path
			if !strings.HasPrefix(expectedPath, "/") && expectedPath != "" {
				expectedPath = "/" + expectedPath
			}

			expectedQuery := test.Query
			if strings.HasPrefix(expectedQuery, "?") && expectedQuery != "" {
				expectedQuery = expectedQuery[1:]
			}

			assert.NotEqual(t, "", short.RedirectPath, "short should have redirect path")
			assert.Equal(t, test.Scheme, short.Scheme, "short scheme shoudl match")
			assert.Equal(t, test.Host, short.Host, "short host should match")
			assert.Equal(t, expectedPath, *short.Path, "short host should match")
			assert.Equal(t, expectedQuery, *short.Query, "short query should match")
			assert.Equal(t, test.RawURL, short.RawURL(), "raw url should be correct")
		})
	}
}
