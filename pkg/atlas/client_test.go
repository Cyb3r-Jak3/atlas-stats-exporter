package atlas

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
	client *API
)

func setup(t *testing.T) {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	client, _ = New(
		WithAPIToken("test-token"),
		WithBaseURL(server.URL),
		WithUserAgent("test-client/1.0"),
	)
}

func Test_ClientOptions(t *testing.T) {
	t.Helper()

	// Create a new API client with test options
	_, err := New(
		WithAPIToken("test-token"),
		WithBaseURL("https://example.com/"),
		WithUserAgent("test-client/1.0"),
		WithHTTPClient(&http.Client{
			Timeout: 30 * time.Second,
		}),
		WithHeaders(http.Header{
			"X-Test-Header": []string{"test-value"},
		},
		),
		WithDebug(true),
	)
	if err != nil {
		t.Fatalf("Failed to create test API client: %v", err)
	}
}

func teardown() {
	server.Close()
}
