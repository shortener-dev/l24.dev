package internal

import (
	"crypto/sha1"
	"encoding/hex"
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

	hasher := sha1.New()
	_, err := hasher.Write([]byte(urlBody))
	if err != nil {
		return nil, err
	}
	hash := hex.EncodeToString(hasher.Sum(nil))

	short := &Short{
		RedirectPath: hash[:7],
		Scheme:       scheme,
		Host:         host,
		Path:         path,
		Query:        query,
	}

	return short, nil
}
