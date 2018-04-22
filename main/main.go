package main

import (
	"fmt"
	"net/http"
	"gophercises/urlshort"
	"log"
	"flag"
	"io/ioutil"
)
const DEFAULT_YAML = "./urls.yaml"
const DEFAULT_JSON = "./urls.json"
var (
	yamlPath = flag.String("yamlPath", DEFAULT_YAML , "a path")
	jsonPath = flag.String("jsonPath", DEFAULT_JSON , "a path")
	)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
		"/fbleite": "http://www.feliperbleite.com",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlBytes, err := ioutil.ReadFile(*yamlPath)
	if err != nil {
		panic(err)
	}
	yamlHandler, err := urlshort.YAMLHandler(yamlBytes, mapHandler)
	if err != nil {
		panic(err)
	}


	jsonBytes, err := ioutil.ReadFile(*jsonPath)
	if err != nil {
		panic(err)
	}
	jsonHandler, err := urlshort.JSONHandler(jsonBytes, yamlHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	err = http.ListenAndServe(":8080", jsonHandler)
	log.Fatal(err)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
