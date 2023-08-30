// Contains methods for converting a tree structure represented by a map of strings to slice of strings.
// The key represents the current parent and the slice of strings represents the childen.
package treeformat

// Used to hold the state during and after conversion.
// Often times this means holding the state of a recursive function.
// It is recommended to use one object per conversion to avoid errors.
type Tree struct {
	// Represents a parent-child relationship between the key and a slice of values.
	// TODO: Make the types more generic in the future.
	relations map[string][]string

	// The string representation of the tree of relations.
	// For example, if we convert the tree to an html list items heirarchy, this will hold the html string.
	s string
}

func (t *Tree) ConstructHtmlList(node string) {
	t.s += "<ul>"
	t.constructHtmlList(node)
	t.s += "</ul>"
}

// Constructs Html unordered list from a map of relations.
// node parameter specifies the root element for the list construction.
func (t *Tree) constructHtmlList(node string) {
	t.s += "<li>" + node

	if _, ok := t.relations[node]; ok && len(t.relations[node]) > 0 {
		// there are children
		t.s += "<ul>"
		for _, link := range t.relations[node] {
			t.constructHtmlList(link)
		}
		t.s += "</ul>"
		t.s += "</li>"
	} else {
		t.s += "</li>"
	}
}
