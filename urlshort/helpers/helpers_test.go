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

// Write some table-driven tests. We can have the name of the test in the struct which is a row in the table.
// Try with subtests too
// check about go 2 error handling 
// as a whole, check the Go 2 draft