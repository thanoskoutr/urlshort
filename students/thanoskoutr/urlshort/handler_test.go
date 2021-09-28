package urlshort

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Testcases Maphandler
var pathsToUrls = map[string]string{
	"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
	"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
}

// Testcases YAMLHandler
var ymls = `
- path: /urlshort-godoc
  url: https://godoc.org/github.com/gophercises/urlshort
- path: /yaml-godoc
  url: https://godoc.org/gopkg.in/yaml.v2
`

func TestMapHandler(t *testing.T) {
	// Run tests for all testcases
	for path, url := range pathsToUrls {
		resp := runMapHandler(t, pathsToUrls, path)

		// Check the status code is what we expect.
		if status := resp.StatusCode; status != http.StatusFound {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusFound)
		}

		// Check the header to see if redirection is what we expect.
		if header := resp.Header; header["Location"][0] != url {
			t.Errorf("handler returned wrong url: got %v want %v",
				header["Location"][0], url)
		}
	}
}

func TestYAMLHandler(t *testing.T) {
	// Run tests for all testcases
	// for i, path := range paths {
	for path, url := range pathsToUrls {
		resp := runYAMLHandler(t, []byte(ymls), path)

		// Check the status code is what we expect.
		if status := resp.StatusCode; status != http.StatusFound {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusFound)
		}

		// Check the header to see if redirection is what we expect.
		if header := resp.Header; header["Location"][0] != url {
			t.Errorf("handler returned wrong url: got %v want %v",
				header["Location"][0], url)
		}
	}

}

// Create a fallback Handler to pass to other Handlers
func fallback(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "fallback handler")
}

// Run the mapHandler with the given path
func runMapHandler(t *testing.T, pathsToUrls map[string]string, path string) *http.Response {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Fatal(err)
	}
	// We create a ResponseRecorder to record the response.
	resp := httptest.NewRecorder()

	// Create a fallback handler
	fallbackHandler := http.HandlerFunc(fallback)

	// Call the handler, passing the response and the request
	mapHandler := MapHandler(pathsToUrls, fallbackHandler)
	mapHandler(resp, req)

	// Return Response: StatusCode, Header, Body
	return resp.Result()
}

// Run the YAMLHandler with the given path
func runYAMLHandler(t *testing.T, yml []byte, path string) *http.Response {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Fatal(err)
	}
	// We create a ResponseRecorder to record the response.
	resp := httptest.NewRecorder()

	// Create a fallback handler
	fallbackHandler := http.HandlerFunc(fallback)

	// Call the handler, passing the response and the request
	yamlHandler, err := YAMLHandler(yml, fallbackHandler)
	if err != nil {
		t.Fatal(err)
	}
	yamlHandler(resp, req)

	// Return Response: StatusCode, Header, Body
	return resp.Result()
}
