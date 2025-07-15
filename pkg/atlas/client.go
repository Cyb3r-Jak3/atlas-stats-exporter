package atlas

import (
	"atlas-stats-exporter/pkg/version"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"regexp"
)

var (
	ErrMissingToken = errors.New("missing API token. Please set the API token using WithAPIToken() option")
)

// API holds the configuration for the current API client. A client should not
// be modified concurrently.
type API struct {
	APIToken   string
	BaseURL    string
	UserAgent  string
	headers    http.Header
	httpClient *http.Client
	Debug      bool
}

func New(opts ...Option) (*API, error) {
	api := &API{
		BaseURL:    "https://atlas.ripe.net/api/v2",
		UserAgent:  fmt.Sprintf("Cyb3rJak3-Atlas-API/%s", version.Version),
		headers:    make(http.Header),
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		if err := opt(api); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return api, nil
}

// copyHeader copies all headers for `source` and sets them on `target`.
// based on https://godoc.org/github.com/golang/gddo/httputil/header#Copy
func copyHeader(target, source http.Header) {
	for k, vs := range source {
		target[k] = vs
	}
}

type APIResponseInfo struct {
	Body       []byte
	StatusCode int
	Headers    http.Header
}

// request makes an HTTP request to the given API endpoint, returning the raw
// *http.Response, or an error if one occurred. The caller is responsible for
// closing the response body.
func (api *API) request(ctx context.Context, method, uri string, reqBody io.Reader, headers http.Header) (*APIResponseInfo, error) {
	req, err := http.NewRequestWithContext(ctx, method, api.BaseURL+uri, reqBody)
	if err != nil {
		return nil, fmt.Errorf("HTTP request creation failed: %w", err)
	}

	combinedHeaders := make(http.Header)
	copyHeader(combinedHeaders, api.headers)
	copyHeader(combinedHeaders, headers)
	req.Header = combinedHeaders

	if api.APIToken != "" {
		req.Header.Set("Authorization", "Key "+api.APIToken)
	}

	if api.UserAgent != "" {
		req.Header.Set("User-Agent", api.UserAgent)
	}

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	if api.Debug {
		dump, httpDumpErr := httputil.DumpRequestOut(req, true)
		if httpDumpErr != nil {
			return nil, httpDumpErr
		}

		// Strip out any sensitive information from the request payload.
		sensitiveKeys := []string{api.APIToken}
		for _, key := range sensitiveKeys {
			if key != "" {
				valueRegex := regexp.MustCompile(fmt.Sprintf("(?m)%s", key))
				dump = valueRegex.ReplaceAll(dump, []byte("[redacted]"))
			}
		}
		log.Printf("\n%s", string(dump))
	}

	resp, err := api.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	if api.Debug {
		dump, httpDumpErr := httputil.DumpResponse(resp, true)
		if httpDumpErr != nil {
			return nil, httpDumpErr
		}
		log.Printf("\n%s", string(dump))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}

	closeErr := resp.Body.Close()
	if closeErr != nil {
		log.Printf("error closing response body: %v", closeErr)
	}

	return &APIResponseInfo{
		Body:       respBody,
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
	}, nil
}
