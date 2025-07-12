package atlas

import (
	"fmt"
	"net/http"
	"time"
)

// Option is a functional option for configuring the API client.
type Option func(*API) error

// WithAPIToken sets the API token for authentication with the Atlas API.
func WithAPIToken(token string) Option {
	return func(api *API) error {
		if token == "" {
			return fmt.Errorf("API token cannot be empty")
		}
		api.APIToken = token
		return nil
	}
}

// WithHTTPClient accepts a custom *http.Client for making API calls.
func WithHTTPClient(client *http.Client) Option {
	return func(api *API) error {
		api.httpClient = client
		return nil
	}
}

// WithHeaders allows you to set custom HTTP headers when making API calls (e.g. for
// satisfying HTTP proxies, or for debugging).
func WithHeaders(headers http.Header) Option {
	return func(api *API) error {
		api.headers = headers
		return nil
	}
}

// WithUserAgent can be set if you want to send a software name and version for HTTP access logs.
// It is recommended to set it in order to help future Customer Support diagnostics
// and prevent collateral damage by sharing generic User-Agent string with abusive users.
func WithUserAgent(userAgent string) Option {
	return func(api *API) error {
		api.UserAgent = userAgent
		return nil
	}
}

// WithBaseURL allows you to override the default HTTP base URL used for API calls.
// This is useful for testing or if you want to use a different API endpoint.
func WithBaseURL(baseURL string) Option {
	return func(api *API) error {
		api.BaseURL = baseURL
		return nil
	}
}

// WithDebug enables or disables debug mode for the API client.
func WithDebug(debug bool) Option {
	return func(api *API) error {
		api.Debug = debug
		return nil
	}
}

// WithTimeout allows you to set a custom timeout for the HTTP client.
func WithTimeout(timeout int) Option {
	return func(api *API) error {
		if timeout <= 0 {
			return fmt.Errorf("timeout must be greater than 0")
		}
		api.httpClient.Timeout = time.Duration(timeout) * time.Second
		return nil
	}
}
