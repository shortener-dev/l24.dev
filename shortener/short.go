package shortener

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
	"time"
)

type Short struct {
	RedirectPath string  `json:"redirect_path" db:"redirect_path"`
	Scheme       string  `json:"scheme" db:"scheme"`
	Host         string  `json:"host" db:"host"`
	Path         *string `json:"path" db:"path"`
	Query        *string `json:"query" db:"query"`
	Fragment     *string `json:"fragment" db:"fragment"`
}

func (s *Short) RawURL() string {
	url := s.Scheme + "://" + s.Host

	if !isNilOrEmptyString(s.Path) {
		url += *s.Path
	}

	if !isNilOrEmptyString(s.Query) {
		url += "?" + *s.Query
	}

	if !isNilOrEmptyString(s.Fragment) {
		url += "#" + *s.Fragment
	}

	return url
}

func NewShort(scheme, host, path, query, fragment string) (*Short, error) {
	if !strings.HasPrefix(path, "/") && path != "" {
		path = "/" + path
	}

	query = strings.TrimPrefix(query, "?")

	fragment = strings.TrimPrefix(fragment, "#")

	short := &Short{
		Scheme:   scheme,
		Host:     host,
		Path:     &path,
		Query:    &query,
		Fragment: &fragment,
	}

	hash, err := createHash(short.RawURL())
	if err != nil {
		return nil, err
	}

	short.RedirectPath = *hash

	return short, nil
}

func createHash(text string) (*string, error) {
	now := time.Now().String() // Dynamic salt to ensure hash uniqueness

	hasher := sha1.New()
	_, err := hasher.Write([]byte(text + now))
	if err != nil {
		return nil, err
	}
	hash := hex.EncodeToString(hasher.Sum(nil))
	hash = hash[:7] // use only first 6 characters - needs to be a short url right?

	return &hash, nil
}
