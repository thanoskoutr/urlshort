package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/thanoskoutr/urlshort/students/thanoskoutr/urlshort"
)

func main() {
	// Parse command-line flag
	yamlFilename := flag.String("yaml", "urls.yaml", "YAML file with URLs and their short paths")
	jsonFilename := flag.String("json", "urls.json", "JSON file with URLs and their short paths")
	flag.Parse()

	// Create a default request multiplexer as the last fallback
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	mapHandler := createMapHandler(mux)

	// Build the YAMLHandler using the previous handler as the fallback
	yamlHandler := createYAMLHandler(*yamlFilename, mapHandler)

	// Build the JSONHandler using the previous handler as the fallback
	jsonHandler := createJSONHandler(*jsonFilename, yamlHandler)

	// Start server
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

// defaultMux is a default request multiplexer for all paths
func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

// hello is a function handler for all paths
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

// createMapHandler creates and returns a Map Hundler
func createMapHandler(fallback http.Handler) http.HandlerFunc {
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, fallback)
	return mapHandler
}

// createYAMLHandler reads the YAML file, creates and returns a YAML Hundler
func createYAMLHandler(name string, fallback http.Handler) http.HandlerFunc {
	yamlFile, err := os.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	yamlHandler, err := urlshort.YAMLHandler([]byte(yamlFile), fallback)
	if err != nil {
		log.Fatal(err)
	}
	return yamlHandler
}

// createJSONHandler reads the JSON file, creates and returns a JSON Hundler
func createJSONHandler(name string, fallback http.Handler) http.HandlerFunc {
	jsonFile, err := os.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	jsonHandler, err := urlshort.JSONHandler([]byte(jsonFile), fallback)
	if err != nil {
		log.Fatal(err)
	}
	return jsonHandler
}
