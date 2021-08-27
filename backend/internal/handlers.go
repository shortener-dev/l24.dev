package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
)

type CreateShortRequest struct {
	URL string `json:"url"`
}

type CreateShortResponse Short

// DecodeJSONBody unmarshalls a JSON response into a struct, while returning any bad request errors
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	contentType := r.Header.Get("Content-Type")
	if contentType != "" {
		if contentType != "application/json" {
			return fmt.Errorf("Content-Type header is '%s', not 'application/json'", contentType)
		}
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&dst)
	return err
}

func NewCreateShortHandler(dao ShortDAO) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request CreateShortRequest

		err := DecodeJSONBody(w, r, &request)
		if err != nil {
			log.Printf("failed to decode json body: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if !strings.HasPrefix(request.URL, "http://") && !strings.HasPrefix(request.URL, "https://") {
			request.URL = "http://" + request.URL // default to http
		}

		URL, err := url.Parse(request.URL)
		if err != nil {
			log.Printf("cannot parse url: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if URL.Host == "" {
			log.Printf("invalid url: %s", URL.String())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		short, err := NewShort(URL.Scheme, URL.Host, URL.Path, URL.RawQuery)
		if err != nil {
			log.Printf("failed to create short: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = dao.InsertShort(*short)
		if err != nil {
			log.Printf("failed to insert short value %s: %v", short.RedirectPath, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(short)
	}
}

func NewGetShortHandler(dao ShortDAO) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		short_url := vars["short"]

		short, err := dao.GetShort(short_url)
		if err != nil {
			log.Printf("failed to find short: %v", err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		http.Redirect(w, r, short.RawURL(), http.StatusMovedPermanently)
	}
}
