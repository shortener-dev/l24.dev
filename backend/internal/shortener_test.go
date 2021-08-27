package internal_test

import (
	"testing"
	"url-shortnener/internal"
)

func TestShortenerFunction(t *testing.T) {
	short, err := internal.NewShort("http", "lucastephens.com", "resume.pdf", "")
	if err != nil {
		t.Fatalf("failed to create short: %v", err)
	}

	if short.RedirectPath == "" {
		t.Fatal("short redirect is empty")
	}

	if short.Host != "lucastephens.com" {
		t.Fatal("host is not lucastephens.com")
	}
}
