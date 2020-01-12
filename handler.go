package urlshort

import (
	"gopkg.in/yaml.v2"
	"net/http"
	"os"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	var hf http.HandlerFunc
	hf = func(w http.ResponseWriter, r *http.Request) {
		u, exist := pathsToUrls[r.URL.Path]
		if exist {
			http.Redirect(w, r, u, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
	return hf
}

type PE struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func parseYAML(yml string) ([]PE, error) {
	file, err := os.Open(yml)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var pes []PE
	err = decoder.Decode(&pes)
	return pes, err
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
func YAMLHandler(yml string, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := make(map[string]string, len(parsedYaml))

	for _, e := range parsedYaml {
		pathMap[e.Path] = e.Url
	}
	return MapHandler(pathMap, fallback), nil
}
