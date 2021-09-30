package urlshort

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/thanoskoutr/urlshort/students/thanoskoutr/database"
	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortPath := r.URL.Path
		if url, ok := pathsToUrls[shortPath]; ok {
			http.Redirect(w, r, url, http.StatusFound)

		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// pathUrl represents the schema of the YAML file, containing paths and their URLs.
type pathUrl struct {
	Url  string
	Path string
}

// parseEncoding will parse an encoded file to validate it.
func parseEncoded(data []byte, enc string) ([]pathUrl, error) {
	var pathUrls []pathUrl
	var err error
	// Select encoding scheme
	switch enc {
	case "yaml":
		err = yaml.Unmarshal(data, &pathUrls)
	case "json":
		err = json.Unmarshal(data, &pathUrls)
	default:
		return nil, fmt.Errorf("%s encoding not supported", enc)
	}
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

// buildMap will convert the parsed data in a YAML file to map.
func buildMap(pathUrls []pathUrl) map[string]string {
	pathUrlMap := make(map[string]string)
	for _, pathUrlItem := range pathUrls {
		pathUrlMap[pathUrlItem.Path] = pathUrlItem.Url
	}
	return pathUrlMap
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYAML, err := parseEncoded(yml, "yaml")
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYAML)
	return MapHandler(pathMap, fallback), nil
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
//    [
//      {
//        "url": "https://www.some-url.com/demo",
//        "path": "/some-path"
//      },
//    ]
//
// The only errors that can be returned all related to having
// invalid JSON data.
func JSONHandler(jsonBlob []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJSON, err := parseEncoded(jsonBlob, "json")
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedJSON)
	return MapHandler(pathMap, fallback), nil
}

// DBHandler will query the Database and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the Database, then the
// fallback http.Handler will be called instead.
//
// Database is expected to be in key-value pair format.
//
// The only errors that can be returned all related to getting
// error from the Database.
func DBHandler(db *database.Database, fallback http.Handler) (http.HandlerFunc, error) {
	pathMap := make(map[string]string)
	// Read all entries from Database, save in a map
	entries, err := database.GetEntriesDB(db)
	if err != nil {
		return nil, err
	}
	for k, v := range entries {
		pathMap[k] = v
	}
	return MapHandler(pathMap, fallback), nil
}
