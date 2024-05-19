package poynt

type NextPageResult struct {
	NextPage     bool
	NextPageLink string
}

func hasNextPage(links []Link) NextPageResult {
	res := NextPageResult{NextPage: false, NextPageLink: ""}
	for _, link := range links {
		if link.Rel == "next" {
			res.NextPage = true
			res.NextPageLink = link.Href
			break
		}
	}
	return res
}
