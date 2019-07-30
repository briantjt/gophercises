package urlshort

import (
	"fmt"
	"net/http"

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
		if redirectURL, ok := pathsToUrls[r.URL.String()]; ok {
			http.Redirect(w, r, redirectURL, 301)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

type urlMap struct {
	Path string
	URL  string
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
	var urlMaps []urlMap
	err := yaml.Unmarshal(yml, &urlMaps)
	if err != nil {
		return nil, err
	}
	return func(w http.ResponseWriter, r *http.Request) {
		original_path := r.URL.String()
		fmt.Printf("Test hello\n")
		for _, urls := range urlMaps {
			if urls.Path == original_path {
				fmt.Printf("Original path is %s\n", original_path)
				http.Redirect(w, r, urls.URL, 301)
				return
			}
		}
		fallback.ServeHTTP(w, r)
	}, nil
}
