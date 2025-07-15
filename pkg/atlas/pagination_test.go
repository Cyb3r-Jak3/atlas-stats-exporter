package atlas

import "testing"

func Test_returnsEmptyPaginationWhenNoMorePages(t *testing.T) {
	p := Pagination{Count: 10, Page: 2, PerPage: 5}
	next := p.Next()
	expected := Pagination{Count: 10, Page: 3, PerPage: 5}
	if next != expected {
		t.Fatalf("expected %+v, got %+v", expected, next)
	}
}

func Test_returnsNextPaginationWhenMorePagesExist(t *testing.T) {
	p := Pagination{Count: 10, Page: 1, PerPage: 5}
	next := p.Next()
	expected := Pagination{Count: 10, Page: 2, PerPage: 5}
	if next != expected {
		t.Fatalf("expected %+v, got %+v", expected, next)
	}
}

func Test_returnsDoneTrueWhenNoMorePages(t *testing.T) {
	p := Pagination{Count: 10, Page: 2, PerPage: 5}
	if !p.Done() {
		t.Fatalf("expected Done to return true, got false")
	}
}

func Test_returnsDoneFalseWhenMorePagesExist(t *testing.T) {
	p := Pagination{Count: 10, Page: 1, PerPage: 5}
	if p.Done() {
		t.Fatalf("expected Done to return false, got true")
	}
}

func Test_URIWithPathAndQueryParameters(t *testing.T) {
	options := Pagination{Count: 10, Page: 1, PerPage: 5}
	expected := "/example-path?page=1&per_page=5"
	result := buildURI("/example-path", options)
	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func Test_returnsDoneTrueWhenNoNext(t *testing.T) {
	p := Pagination{Count: 10, Page: 2, PerPage: 5}
	if !p.Done() {
		t.Fatalf("expected Done to return true, got false")
	}
}
