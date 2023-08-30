// Contains methods for converting a tree structure represented by a map of strings to slice of strings.
// The key represents the current parent and the slice of strings represents the childen.
package treeformat

// Used to hold the state during and after conversion.
// Often times this means holding the state of a recursive function.
// It is recommended to use one object per conversion to avoid errors.
type Tree struct {
	// Represents a parent-child relationship between the key and a slice of values.
	// The client has to set this field before calling methods on this object.
	Relations map[string][]string

	// The string representation of the tree of relations.
	// For example, if we convert the tree to an html list items heirarchy, this will hold the html string.
	S string

	// If we already visited a link we might list it as a child, but will not recurse in it.
	visited map[string]bool
}

func (t *Tree) ConstructHtmlList(node string) {
	t.visited = make(map[string]bool)

	t.S += "<ul>"
	t.constructHtmlList(node)
	t.S += "</ul>"
}

// Constructs Html unordered list from a map of relations.
// node parameter specifies the root element for the list construction.
func (t *Tree) constructHtmlList(node string) {
	t.S += "<li>" + node

	_, isVisited := t.visited[node]
	if _, ok := t.Relations[node]; ok && len(t.Relations[node]) > 0 && !isVisited {
		// there are children

		t.visited[node] = true
		t.S += "<ul>"

		for _, link := range t.Relations[node] {
			t.constructHtmlList(link)
		}
		t.S += "</ul>"
		t.S += "</li>"
	} else {
		t.S += "</li>"
	}
}
