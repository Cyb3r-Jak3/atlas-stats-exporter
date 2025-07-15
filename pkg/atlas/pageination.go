package atlas

import (
	"log"
	"net/url"

	"github.com/google/go-querystring/query"
)

type Pagination struct {
	Count    int    `json:"count" url:"-"`
	Page     int    `url:"page,omitempty"`
	PerPage  int    `url:"per_page,omitempty"`
	NextLink string `json:"next,omitempty" url:"-"`
}

func (p Pagination) Next() Pagination {
	return Pagination{
		Count:    p.Count,
		Page:     p.Page + 1,
		PerPage:  p.PerPage,
		NextLink: p.NextLink,
	}
}

func (p Pagination) Done() bool {
	//noNext := p.NextLink == ""
	if p.NextLink != "" {
		return false
	}
	log.Printf("Pagination: Count=%d, Page=%d, PerPage=%d, NextLink=%s\n", p.Count, p.Page, p.PerPage, p.NextLink)
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
