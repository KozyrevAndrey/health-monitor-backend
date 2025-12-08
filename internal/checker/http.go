package checker

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"

	"health-monitor/internal/domain"
)

type HTTPChecker struct {
	client *http.Client
}

func NewHTTPChecker() *HTTPChecker {
	return &HTTPChecker{
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				MaxIdleConns:          100,
				MaxIdleConnsPerHost:   10,
				IdleConnTimeout:       90 * time.Second,
			},
		},
	}
}

func (c *HTTPChecker) Check(ctx context.Context, target *domain.Target) (*domain.CheckResult, error) {
	httpConfig, err := parseHTTPConfig(target.Config)
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP config: %w", err)
	}

	result := &domain.CheckResult{
		TargetID:  target.ID,
		CheckedAt: time.Now(),
		Metadata:  make(map[string]interface{}),
	}

	checkCtx, cancel := context.WithTimeout(ctx, target.Timeout)
	defer cancel()

	startTime := time.Now()
	statusCode, err := c.performHTTPRequest(checkCtx, httpConfig, result)
	responseTime := time.Since(startTime)

	result.ResponseTimeMs = responseTime.Milliseconds()
	result.StatusCode = statusCode

	if err != nil {
		result.Status = domain.CheckStatusFailure
		result.Error = err.Error()
		result.Message = fmt.Sprintf("HTTP check failed: %s", err.Error())
		return result, nil
	}

	if statusCode != httpConfig.ExpectedStatusCode {
		result.Status = domain.CheckStatusFailure
		result.Message = fmt.Sprintf("Unexpected status code: got %d, expected %d", statusCode, httpConfig.ExpectedStatusCode)
		return result, nil
	}

	if httpConfig.MaxResponseTimeMs > 0 && responseTime.Milliseconds() > int64(httpConfig.MaxResponseTimeMs) {
		result.Status = domain.CheckStatusWarning
		result.Message = fmt.Sprintf("Slow response: %dms (threshold: %dms)", responseTime.Milliseconds(), httpConfig.MaxResponseTimeMs)
		return result, nil
	}

	result.Status = domain.CheckStatusSuccess
	result.Message = fmt.Sprintf("HTTP check successful: %d %s in %dms", statusCode, http.StatusText(statusCode), responseTime.Milliseconds())

	return result, nil
}

func (c *HTTPChecker) Type() domain.TargetType {
	return domain.TargetTypeHTTP
}

func (c *HTTPChecker) Validate(config map[string]interface{}) error {
	_, err := parseHTTPConfig(config)
	return err
}

func (c *HTTPChecker) performHTTPRequest(ctx context.Context, config *domain.HTTPConfig, result *domain.CheckResult) (int, error) {
	req, err := http.NewRequestWithContext(ctx, config.Method, config.URL, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "HealthMonitor/1.0")
	}

	client := c.client
	if config.FollowRedirects == false || config.ValidateSSL == false {
		transport := c.client.Transport.(*http.Transport).Clone()

		if config.ValidateSSL == false {
			transport.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		client = &http.Client{
			Timeout:   c.client.Timeout,
			Transport: transport,
		}

		if config.FollowRedirects == false {
			client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	io.Copy(io.Discard, resp.Body)

	result.Metadata["url"] = config.URL
	result.Metadata["method"] = config.Method
	result.Metadata["status_text"] = resp.Status

	if config.CheckSSLExpiry && resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
		cert := resp.TLS.PeerCertificates[0]
		daysUntilExpiry := int(time.Until(cert.NotAfter).Hours() / 24)
		result.Metadata["ssl_expiry_days"] = daysUntilExpiry
		result.Metadata["ssl_expires_at"] = cert.NotAfter.Format(time.RFC3339)

		if config.SSLExpiryDays > 0 && daysUntilExpiry <= config.SSLExpiryDays {
			result.Metadata["ssl_warning"] = fmt.Sprintf("SSL certificate expires in %d days", daysUntilExpiry)
		}
	}

	return resp.StatusCode, nil
}

func parseHTTPConfig(config map[string]interface{}) (*domain.HTTPConfig, error) {
	httpConfig := &domain.HTTPConfig{
		Method:             "GET",
		ExpectedStatusCode: 200,
		FollowRedirects:    true,
		ValidateSSL:        true,
	}

	if url, ok := config["url"].(string); ok && url != "" {
		httpConfig.URL = url
	} else {
		return nil, fmt.Errorf("url is required")
	}

	if method, ok := config["method"].(string); ok && method != "" {
		httpConfig.Method = method
	}

	if headers, ok := config["headers"].(map[string]interface{}); ok {
		httpConfig.Headers = make(map[string]string)
		for key, value := range headers {
			if strValue, ok := value.(string); ok {
				httpConfig.Headers[key] = strValue
			}
		}
	}

	if body, ok := config["body"].(string); ok {
		httpConfig.Body = body
	}

	if statusCode, ok := config["expected_status_code"].(int); ok {
		httpConfig.ExpectedStatusCode = statusCode
	} else if statusCode, ok := config["expected_status_code"].(float64); ok {
		httpConfig.ExpectedStatusCode = int(statusCode)
	}

	if maxTime, ok := config["max_response_time_ms"].(int); ok {
		httpConfig.MaxResponseTimeMs = maxTime
	} else if maxTime, ok := config["max_response_time_ms"].(float64); ok {
		httpConfig.MaxResponseTimeMs = int(maxTime)
	}

	if follow, ok := config["follow_redirects"].(bool); ok {
		httpConfig.FollowRedirects = follow
	}

	if validate, ok := config["validate_ssl"].(bool); ok {
		httpConfig.ValidateSSL = validate
	}

	if check, ok := config["check_ssl_expiry"].(bool); ok {
		httpConfig.CheckSSLExpiry = check
	}

	if days, ok := config["ssl_expiry_days"].(int); ok {
		httpConfig.SSLExpiryDays = days
	} else if days, ok := config["ssl_expiry_days"].(float64); ok {
		httpConfig.SSLExpiryDays = int(days)
	}

	return httpConfig, nil
}
