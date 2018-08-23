package urlshort

import (
	"fmt"
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// if we can match a path... redirect to it
		// if their is a vaile ok will be true, otherwise false.
		// The dest variable will contain a value if there is a match
		if dest, ok := pathsToUrls[path]; ok {
			fmt.Printf("dest, ok: %s, %t\n", string(dest), ok)
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		// otherwise
		fallback.ServeHTTP(w, r)
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
func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// 1. Parse the YAML somehow
	pathURLs, err := parseYaml(yamlBytes)
	if err != nil {
		return nil, err
	}
	// 2. convert YAML array into map
	pathsToUrls := buildMap(pathURLs)
	// 3. return a map handler using the mapping
	return MapHandler(pathsToUrls, fallback), nil

}

func parseYaml(data []byte) ([]pathURL, error) {
	var pathURLs []pathURL
	err := yaml.Unmarshal(data, &pathURLs)
	if err != nil {
		return nil, err
	}
	return pathURLs, nil
}

func buildMap(pathURLs []pathURL) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, pu := range pathURLs {
		pathsToUrls[pu.Path] = pu.URL
	}
	return pathsToUrls
}

// TODO try these members lower case
type pathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}
