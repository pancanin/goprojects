package helpers

import "fmt"

type PathUrl struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func CreateUrlToPathMap(pathUrls []PathUrl) map[string]string {
	pathToUrlMap := make(map[string]string, len(pathUrls))

	for _, pathUrl := range pathUrls {
		pathToUrlMap[pathUrl.Path] = pathUrl.Url
		fmt.Printf("The value added is %s and the key is %s", pathUrl.Url, pathUrl.Path)
	}

	return pathToUrlMap
}
