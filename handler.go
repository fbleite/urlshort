package urlshort

import (
	"net/http"
	"log"
	"gopkg.in/yaml.v2"
	"encoding/json"
	"github.com/boltdb/bolt"
	"fmt"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request path", r.URL.Path)
		if redirectUrl, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, redirectUrl, 301)
		} else {
			fallback.ServeHTTP(w, r)
		}
		return
	}
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
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMapYaml(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func parseYAML(yamlInput []byte) ([]pathUrlYaml, error) {
	var parsedYaml []pathUrlYaml
	err :=yaml.Unmarshal(yamlInput, &parsedYaml)
	return parsedYaml,err
}

type pathUrlYaml struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func buildMapYaml(parsedYaml []pathUrlYaml) map[string]string {
	pathUrlMap := make(map[string]string)
	for _, tuple := range(parsedYaml) {
		pathUrlMap[tuple.Path] = tuple.URL
	}
	return pathUrlMap
}


func JSONHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJson, err := parseJSON(json)
	if err != nil {
		return nil, err
	}
	urlMap := buildMapJson(parsedJson)
	return MapHandler(urlMap, fallback), nil
}

func parseJSON(jsonInput []byte) ([]pathUrlJson, error) {
	var parsedJson []pathUrlJson
	err :=json.Unmarshal(jsonInput, &parsedJson)
	return parsedJson,err
}


type pathUrlJson struct {
	Path string `json:path`
	URL  string `json:url`
}


func buildMapJson(parsedJson []pathUrlJson) map[string]string {
	pathUrlMap := make(map[string]string)
	for _, tuple := range(parsedJson) {
		pathUrlMap[tuple.Path] = tuple.URL
	}
	return pathUrlMap
}

func BoltHandler(db bolt.DB, fallback http.Handler) (http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request path", r.URL.Path)
		var url string
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("path-url"))
			url = string(b.Get([]byte(r.URL.Path)))
			fmt.Printf("The answer is: %s\n", url)
			return nil
		})

		if url != "" {
			http.Redirect(w, r, url, 301)
		} else {
			fallback.ServeHTTP(w, r)
		}
		return
	}

}

