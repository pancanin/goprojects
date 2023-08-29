package treeformat

import (
	"reflect"
	"testing"
)

func TestConstructHtmlList(t *testing.T) {
	tests := map[string]struct {
		input map[string][]string
		want string
	}{
		"flat list": {
			input: map[string][]string{
				"/": {"/articles", "/about", "/contacts"},
			},
			want: "<ul><li>/<ul><li>/articles</li><li>/about</li><li>/contacts</li></ul></li></ul>",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tree := Tree{relations: tc.input}
			tree.ConstructHtmlList("/")

			if !reflect.DeepEqual(tree.s, tc.want) {
				t.Fatalf("expected:\n%v,\n got:\n%v\n", tc.want, tree.s)
			}
		})
	}
}