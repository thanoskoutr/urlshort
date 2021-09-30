package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/thanoskoutr/urlshort/students/thanoskoutr/database"
	"github.com/thanoskoutr/urlshort/students/thanoskoutr/urlshort"
)

func main() {
	// Parse command-line flag
	yamlFilename := flag.String("yaml", "urls.yaml", "YAML file with URLs and their short paths")
	jsonFilename := flag.String("json", "urls.json", "JSON file with URLs and their short paths")
	dbFilename := flag.String("db", "urls.db", "Database file")
	flag.Parse()

	// Setup Database
	BUCKET_NAME := "URL"
	db, err := database.SetupDB(*dbFilename, BUCKET_NAME)
	if err != nil {
		log.Fatal(err)
	}
	defer db.BoltDB.Close()

	// Create a default request multiplexer as the last fallback
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	mapHandler := createMapHandler(mux)

	// Build the YAMLHandler using the previous handler as the fallback
	yamlHandler := createYAMLHandler(*yamlFilename, mapHandler)

	// Build the JSONHandler using the previous handler as the fallback
	jsonHandler := createJSONHandler(*jsonFilename, yamlHandler)

	// Build the DBHandler using the previous handler as the fallback
	dbHandler := createDBHandler(db, jsonHandler)

	// Start server
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", dbHandler)
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

// createDBHandler reads the Database, creates and returns a DB Hundler
func createDBHandler(db *database.Database, fallback http.Handler) http.HandlerFunc {
	// Add some entries to the Database
	pathsToUrls := map[string]string{
		"/gnu/health":   "https://savannah.gnu.org/projects/health",
		"/gnu/avr-libc": "https://savannah.nongnu.org/projects/avr-libc",
		"/gnu/dino":     "https://savannah.nongnu.org/projects/dino",
		"/gnu/ddd":      "https://savannah.gnu.org/projects/ddd",
		"/gnu/epsilon":  "https://savannah.gnu.org/projects/epsilon",
	}
	err := database.PutMapEntriesDB(db, pathsToUrls)
	if err != nil {
		log.Fatal(err)
	}

	dbHandler, err := urlshort.DBHandler(db, fallback)
	if err != nil {
		log.Fatal(err)
	}
	return dbHandler
}
