package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"example.com/urlshortener/urlshort"
)

func main() {
	mux := defaultMux()

	// We have a few hardcoded mappings...The developers were lazy to add them to yaml or json... :D
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// If the mappings are not in the hardcoded map above, check the yaml file.
	// Yaml files have really quirky idented syntax and I don't like them much :D
	yamlFileName := flag.String("yaml", "urls.yaml", "path to yaml file with path, url pairs")
	yamlFileContent := readFileContents(*yamlFileName)
	
	yamlHandler, err := urlshort.YAMLHandler(yamlFileContent, mapHandler)
	if err != nil {
		panic(err)
	}

	jsonFileName := flag.String("json", "urls.json", "path to json formatted plain-text file")
	jsonFileContent := readFileContents(*jsonFileName)

	jsonHandler, err := urlshort.JSONHandler(jsonFileContent, yamlHandler)

	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func readFileContents(path string) []byte {
	contents, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open file %s", path)
		os.Exit(1)
	}
	return contents
}