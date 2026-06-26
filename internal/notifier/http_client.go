package notifier

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// proxySchemes lists the proxy URL schemes supported by buildHTTPClient.
// net/http's Transport understands all of these natively (HTTP/HTTPS CONNECT
// proxies and SOCKS5), so no extra dependency is required.
var proxySchemes = map[string]bool{
	"http":   true,
	"https":  true,
	"socks5": true,
}

// buildHTTPClient returns an *http.Client with the given timeout. When proxyURL
// is non-empty it is parsed and validated, and all requests are routed through
// it. Supported schemes: http, https, socks5.
func buildHTTPClient(proxyURL string, timeout time.Duration) (*http.Client, error) {
	client := &http.Client{Timeout: timeout}

	if proxyURL == "" {
		return client, nil
	}

	parsed, err := parseProxyURL(proxyURL)
	if err != nil {
		return nil, err
	}

	client.Transport = &http.Transport{
		Proxy: http.ProxyURL(parsed),
	}
	return client, nil
}

// parseProxyURL validates a proxy URL and returns the parsed value.
func parseProxyURL(proxyURL string) (*url.URL, error) {
	parsed, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("invalid proxy_url: %w", err)
	}
	if !proxySchemes[parsed.Scheme] {
		return nil, fmt.Errorf("unsupported proxy scheme %q (supported: http, https, socks5)", parsed.Scheme)
	}
	if parsed.Host == "" {
		return nil, fmt.Errorf("proxy_url must include a host")
	}
	return parsed, nil
}
