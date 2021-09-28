package urlshort

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Testcases
var pathsToUrls = map[string]string{
	"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
	"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
}

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

func TestParseYAML(t *testing.T) {
	t.Errorf("not implemented")
}

func TestBuildMap(t *testing.T) {
	t.Errorf("not implemented")
}

func TestYAMLHandler(t *testing.T) {
	t.Errorf("not implemented")
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
