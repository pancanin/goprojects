package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"example.com/urlshortener/handler"
)

func main() {
	mux := defaultMux()

	// Hardcoded mappings.
	pathToUrl := map[string]string{
		"/handler-godoc": "https://godoc.org/github.com/gophercises/handler",
		"/yaml-godoc":    "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := handler.MapHandler(pathToUrl, mux)

	// If the mappings are not in the hardcoded map above, check the yaml file.
	yamlFileName := flag.String("yaml", "urls.yaml", "path to yaml file with path, url pairs")
	yamlFileContent := readFileContents(*yamlFileName)

	yamlHandler, err := handler.YAMLHandler(yamlFileContent, mapHandler)
	if err != nil {
		fmt.Fprintf(os.Stderr, "There was a problem with reading yaml contents: %s", err.Error())
		os.Exit(1)
	}

	jsonFileName := flag.String("json", "urls.json", "path to json formatted plain-text file")
	jsonFileContent := readFileContents(*jsonFileName)
	jsonHandler, err := handler.JSONHandler(jsonFileContent, yamlHandler)

	if err != nil {
		fmt.Fprintf(os.Stderr, "There was a problem with reading json contents: %s", err.Error())
		os.Exit(1)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", defaultHandle)
	return mux
}

func defaultHandle(w http.ResponseWriter, r *http.Request) {
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
