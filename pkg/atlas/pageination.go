package atlas

import (
	"net/url"

	"github.com/google/go-querystring/query"
)

type Pagination struct {
	Count   int `json:"count" url:"count,omitempty"`
	Page    int `json:"page" url:"page,omitempty"`
	PerPage int `json:"per_page" url:"per_page,omitempty"`
}

func (p Pagination) Next() Pagination {
	if p.Count == 0 || p.Page >= (p.Count/p.PerPage) {
		return Pagination{}
	}
	return Pagination{
		Count:   p.Count,
		Page:    p.Page + 1,
		PerPage: p.PerPage,
	}
}

func (p Pagination) Done() bool {
	if p.Count == 0 || p.Page >= (p.Count/p.PerPage) {
		return true
	}
	return false
}

// buildURI assembles the base path and queries.
func buildURI(path string, options interface{}) string {
	v, _ := query.Values(options)
	return (&url.URL{Path: path, RawQuery: v.Encode()}).String()
}
