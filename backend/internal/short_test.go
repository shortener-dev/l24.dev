package internal_test

import (
	"testing"
	"url-shortnener/internal"

	"github.com/stretchr/testify/assert"
)

func TestNewShort(t *testing.T) {
	short1, err := internal.NewShort("http", "lucastephens.com", "resume.pdf", "")
	if err != nil {
		t.Fatalf("failed to create short: %v", err)
	}

	short2, err := internal.NewShort("http", "lucastephens.com", "resume.pdf", "")
	if err != nil {
		t.Fatalf("failed to create short: %v", err)
	}

	assert.NotEqual(t, "", short1.RedirectPath, "short1 should have redirect path")
	assert.NotEqual(t, "", short2.RedirectPath, "short2 should have redirect path")
	assert.Equal(t, "lucastephens.com", short1.Host, "short1 host should match")
	assert.Equal(t, "lucastephens.com", short2.Host, "short2 host should match")
	assert.Equal(t, "resume.pdf", short1.Path, "short1 host should match")
	assert.Equal(t, "resume.pdf", short2.Path, "short2 path should match")
	assert.Equal(t, "", short1.Query, "short1 query should match")
	assert.Equal(t, "", short2.Query, "short2 query should match")
	assert.NotEqual(t, short1.RedirectPath, short2.RedirectPath, "shorts should have different hashes")
}
