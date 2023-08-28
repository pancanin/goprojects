package link

type Link struct {
	Path string
	Url string
}

func TransformLinksToMap(links []Link) map[string]string {
	res := make(map[string]string, len(links))

	for _, l := range links {
		res[l.Path] = l.Url
	}

	return res
}
