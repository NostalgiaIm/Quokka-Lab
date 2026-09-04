package rails

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// Client owns internal communication from the Go gateway to Rails.
type Client struct {
	baseURL *url.URL
	proxy   *httputil.ReverseProxy
	http    *http.Client
	logger  *slog.Logger
}

// NewClient creates an internal Rails client. The current MVP uses HTTP proxying;
// backend/proto defines the planned gRPC contract.
func NewClient(rawURL string, logger *slog.Logger) *Client {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(parsedURL)
	originalDirector := proxy.Director
	proxy.Director = func(r *http.Request) {
		originalDirector(r)
		r.Host = parsedURL.Host
		r.Header.Set("X-Quokka-Internal-Caller", "go-gateway")
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		logger.Error("rails proxy failed", "path", r.URL.Path, "error", err)
		http.Error(w, "rails service unavailable", http.StatusBadGateway)
	}

	return &Client{
		baseURL: parsedURL,
		proxy:   proxy,
		http:    &http.Client{Timeout: 90 * time.Second},
		logger:  logger,
	}
}

// Proxy forwards the request to Rails while preserving method, body, and headers.
func (c *Client) Proxy(w http.ResponseWriter, r *http.Request) {
	c.proxy.ServeHTTP(w, r)
}

// ForwardMultipart sends an upload request to Rails and decodes the JSON response.
func (c *Client) ForwardMultipart(r *http.Request, targetPath string) (map[string]any, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	target := *c.baseURL
	target.Path = targetPath
	target.RawQuery = r.URL.RawQuery

	request, err := http.NewRequestWithContext(r.Context(), r.Method, target.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	copyHeaders(request.Header, r.Header)
	request.Header.Set("X-Quokka-Internal-Caller", "go-gateway")

	response, err := c.http.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	payload, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode >= 400 {
		return nil, railsError{status: response.StatusCode, body: strings.TrimSpace(string(payload))}
	}

	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		return nil, err
	}
	return decoded, nil
}

func copyHeaders(dst http.Header, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

type railsError struct {
	status int
	body   string
}

func (e railsError) Error() string {
	if e.body == "" {
		return http.StatusText(e.status)
	}
	return e.body
}
