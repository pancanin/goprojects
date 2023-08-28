package link_test

import (
	"reflect"
	"testing"

	"example.com/urlshortener/link"
)

func TestCreateUrlToPathMap(t *testing.T) {
	tests := map[string]struct {
		input []link.Link
		want map[string]string
	}{
		"simple": {
			input: []link.Link{
				link.Link{Path: "/a", Url: "www.abc.bg"},
				link.Link{Path: "/b", Url: "www.bavc.bg"},
			},
			want: map[string]string{
				"/a": "www.abc.bg",
				"/b": "www.bavc.bg",
			},
		},
	}


	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := link.TransformLinksToMap(tc.input)

			if !reflect.DeepEqual(tc.want, actual) {
				t.Fatalf("expected: %v, actual: %v", tc.want, actual)
			}
		})
	}
}

// Write some table-driven tests. We can have the name of the test in the struct which is a row in the table.
// Try with subtests too
// check about go 2 error handling
// as a whole, check the Go 2 draft
