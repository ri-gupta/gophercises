package urlmapper

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

type redirectionData struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if redirectionPath, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, redirectionPath, http.StatusFound)
			return
		}

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
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYaml(yml)
	if err != nil {
		return nil, err
	}

	urlMap := buildMap(parsedYaml)

	return MapHandler(urlMap, fallback), nil
}

func parseYaml(yml []byte) ([]redirectionData, error) {
	var parsedData []redirectionData

	err := yaml.Unmarshal(yml, &parsedData)
	if err != nil {
		err = fmt.Errorf("invalid yaml, unable to parse with error: %v", err)
	}

	return parsedData, err
}

func buildMap(parsedYaml []redirectionData) map[string]string {
	var redirectionMap = make(map[string]string)

	for _, data := range parsedYaml {
		redirectionMap[data.Path] = data.Url
	}

	return redirectionMap
}
