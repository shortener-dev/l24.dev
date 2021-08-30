package internal

import (
	"crypto/sha1"
	"encoding/hex"
	"time"
)

type Short struct {
	RedirectPath string `json:"redirect_path" db:"redirect_path"`
	Scheme       string `json:"schema" db:"scheme"`
	Host         string `json:"host" db:"host"`
	Path         string `json:"path" db:"path"`
	Query        string `json:"query" db:"query"`
}

func (s *Short) RawURL() string {
	if s.Query != "" {
		return s.Scheme + "://" + s.Host + s.Path + "?" + s.Query
	}
	return s.Scheme + "://" + s.Host + s.Path
}

func NewShort(scheme, host, path, query string) (*Short, error) {
	var urlBody string

	if query != "" {
		urlBody = host + path + "?" + query
	} else {
		urlBody = host + path
	}

	hash, err := createHash(urlBody)
	if err != nil {
		return nil, err
	}

	short := &Short{
		RedirectPath: *hash,
		Scheme:       scheme,
		Host:         host,
		Path:         path,
		Query:        query,
	}

	return short, nil
}

func createHash(text string) (*string, error) {
	now := time.Now().String() // This ensures the hash is always unique, since time always moves forward

	hasher := sha1.New()
	_, err := hasher.Write([]byte(text + now))
	if err != nil {
		return nil, err
	}
	hash := hex.EncodeToString(hasher.Sum(nil))
	hash = hash[:7] // use only first 6 characters - needs to be a short url right?

	return &hash, nil
}
