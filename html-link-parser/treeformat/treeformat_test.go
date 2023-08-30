package treeformat

import (
	"reflect"
	"testing"
)

func TestConstructHtmlList(t *testing.T) {
	tests := map[string]struct {
		input map[string][]string
		want  string
	}{
		"flat list": {
			input: map[string][]string{
				"/": {"/articles", "/about", "/contacts"},
			},
			want: "<ul><li>/<ul><li>/articles</li><li>/about</li><li>/contacts</li></ul></li></ul>",
		},
		"nested list": {
			input: map[string][]string{
				"/":      {"/articles", "/blog", "/about"},
				"/about": {"/testimonials", "/pricing"},
				"/blog": {"/golang", "/cpp", "/java"},
			},
			want: "<ul><li>/<ul><li>/articles</li><li>/blog<ul><li>/golang</li><li>/cpp</li><li>/java</li></ul></li><li>/about<ul><li>/testimonials</li><li>/pricing</li></ul></li></ul></li></ul>",
		},
		"cyclic references": {
			input: map[string][]string{
				"/": {"/about"},
				"/about": {"/", "/blog"},
			},
			want: "<ul><li>/<ul><li>/about<ul><li>/</li><li>/blog</li></ul></li></ul></li></ul>",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tree := Tree{Relations: tc.input}
			tree.ConstructHtmlList("/")

			if !reflect.DeepEqual(tree.S, tc.want) {
				t.Fatalf("expected:\n%v,\n got:\n%v\n", tc.want, tree.S)
			}
		})
	}
}
