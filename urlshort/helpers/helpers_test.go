package helpers_test

import (
	"testing"

	"example.com/urlshortener/helpers"
)

func TestCreateUrlToPathMap(t *testing.T) {
	pathToUrlList := []helpers.PathUrl{
		helpers.PathUrl{Path: "/a", Url: "www.abc.bg"},
		helpers.PathUrl{Path: "/b", Url: "www.bavc.bg"},
	}

	resMap := helpers.CreateUrlToPathMap(pathToUrlList)

	if _, ok := resMap["/a"]; !ok {
		t.Fatalf("The map does not contain key '/a")
	}

	// A few similar assertions
}