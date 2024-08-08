package urlshort

import (
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

type pathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if url, found := pathsToUrls[r.URL.Path]; found {
			http.Redirect(w, r, url, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)

	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	var pathURLs []pathURL
	err := yaml.Unmarshal(yml, &pathURLs)
	if err != nil {
		log.Fatalf("error while parsing %s", err)
		os.Exit(1)
	}
	pathsToUrls := map[string]string{}

	for _, item := range pathURLs {
		pathsToUrls[item.Path] = item.URL
	}

	return MapHandler(pathsToUrls, fallback), err

}
