package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"

	"github.com/rs/zerolog"
	"health-monitor/internal/domain"
)

const (
	webhookDefaultMethod     = http.MethodPost
	webhookDefaultTimeout    = 10 * time.Second
	webhookDefaultMaxRetries = 3
	webhookRetryBaseDelay    = 500 * time.Millisecond
)

// WebhookNotifier sends alerts to an arbitrary HTTP endpoint. The request body
// is built from a configurable Go text/template (defaulting to the alert as
// JSON), so it can target Slack, Discord, or any custom webhook.
type WebhookNotifier struct {
	url        string
	method     string
	headers    map[string]string
	tmpl       *template.Template // nil => default JSON body
	maxRetries int
	client     *http.Client
	log        zerolog.Logger
}

type webhookConfig struct {
	URL        string
	Method     string
	Headers    map[string]string
	Payload    string
	Timeout    time.Duration
	MaxRetries int
	ProxyURL   string
}

func NewWebhookNotifier(config map[string]interface{}, log zerolog.Logger) (*WebhookNotifier, error) {
	cfg, err := parseWebhookConfig(config)
	if err != nil {
		return nil, err
	}

	var tmpl *template.Template
	if cfg.Payload != "" {
		tmpl, err = template.New("webhook").Parse(cfg.Payload)
		if err != nil {
			return nil, fmt.Errorf("invalid payload template: %w", err)
		}
	}

	client, err := buildHTTPClient(cfg.ProxyURL, cfg.Timeout)
	if err != nil {
		return nil, err
	}

	return &WebhookNotifier{
		url:        cfg.URL,
		method:     cfg.Method,
		headers:    cfg.Headers,
		tmpl:       tmpl,
		maxRetries: cfg.MaxRetries,
		client:     client,
		log:        log,
	}, nil
}

func parseWebhookConfig(config map[string]interface{}) (*webhookConfig, error) {
	rawURL, ok := config["url"].(string)
	if !ok || rawURL == "" {
		return nil, fmt.Errorf("url is required")
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, fmt.Errorf("url must be http or https")
	}
	if parsed.Host == "" {
		return nil, fmt.Errorf("url must include a host")
	}

	cfg := &webhookConfig{
		URL:        rawURL,
		Method:     webhookDefaultMethod,
		Timeout:    webhookDefaultTimeout,
		MaxRetries: webhookDefaultMaxRetries,
	}

	if raw, ok := config["method"]; ok {
		method, ok := raw.(string)
		if !ok {
			return nil, fmt.Errorf("method must be a string")
		}
		if method != "" {
			cfg.Method = strings.ToUpper(method)
		}
	}

	if raw, ok := config["headers"]; ok {
		headers, err := parseWebhookHeaders(raw)
		if err != nil {
			return nil, err
		}
		cfg.Headers = headers
	}

	if raw, ok := config["payload"]; ok {
		payload, ok := raw.(string)
		if !ok {
			return nil, fmt.Errorf("payload must be a string")
		}
		cfg.Payload = payload
	}

	if raw, ok := config["timeout"]; ok {
		timeoutStr, ok := raw.(string)
		if !ok {
			return nil, fmt.Errorf("timeout must be a string (e.g. \"10s\")")
		}
		if timeoutStr != "" {
			d, err := time.ParseDuration(timeoutStr)
			if err != nil {
				return nil, fmt.Errorf("invalid timeout: %w", err)
			}
			if d <= 0 {
				return nil, fmt.Errorf("timeout must be positive")
			}
			cfg.Timeout = d
		}
	}

	if raw, ok := config["max_retries"]; ok {
		n, err := toInt(raw)
		if err != nil {
			return nil, fmt.Errorf("max_retries must be a number")
		}
		if n < 0 {
			return nil, fmt.Errorf("max_retries must be >= 0")
		}
		cfg.MaxRetries = n
	}

	if raw, ok := config["proxy_url"]; ok {
		proxyURL, ok := raw.(string)
		if !ok {
			return nil, fmt.Errorf("proxy_url must be a string")
		}
		if proxyURL != "" {
			if _, err := parseProxyURL(proxyURL); err != nil {
				return nil, err
			}
		}
		cfg.ProxyURL = proxyURL
	}

	return cfg, nil
}

// parseWebhookHeaders accepts a JSON object of string->string headers.
func parseWebhookHeaders(raw interface{}) (map[string]string, error) {
	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("headers must be an object of string values")
	}
	headers := make(map[string]string, len(m))
	for k, v := range m {
		strVal, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("header %q must be a string", k)
		}
		headers[k] = strVal
	}
	return headers, nil
}

// toInt coerces JSON numbers (float64) or int-like values to int.
func toInt(raw interface{}) (int, error) {
	switch v := raw.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	case int64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("not a number")
	}
}

func (w *WebhookNotifier) Notify(ctx context.Context, alert *domain.Alert) error {
	body, err := w.renderPayload(alert)
	if err != nil {
		return fmt.Errorf("failed to render payload: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt <= w.maxRetries; attempt++ {
		if attempt > 0 {
			delay := webhookRetryBaseDelay * time.Duration(1<<(attempt-1))
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}

		retryable, err := w.send(ctx, body)
		if err == nil {
			if attempt > 0 {
				w.log.Info().Str("alert_id", alert.ID).Int("attempt", attempt+1).Msg("Webhook notification sent after retry")
			} else {
				w.log.Info().Str("alert_id", alert.ID).Str("alert_type", string(alert.Type)).Msg("Webhook notification sent")
			}
			return nil
		}

		lastErr = err
		if !retryable {
			break
		}
		w.log.Warn().Err(err).Str("alert_id", alert.ID).Int("attempt", attempt+1).Msg("Webhook attempt failed, will retry")
	}

	w.log.Error().Err(lastErr).Str("alert_id", alert.ID).Msg("Failed to send webhook notification")
	return fmt.Errorf("webhook notification failed: %w", lastErr)
}

// send performs a single request. The bool indicates whether the failure is
// retryable (network errors, 429, 5xx).
func (w *WebhookNotifier) send(ctx context.Context, body []byte) (retryable bool, err error) {
	req, err := http.NewRequestWithContext(ctx, w.method, w.url, bytes.NewReader(body))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	if _, hasContentType := w.headers["Content-Type"]; !hasContentType {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range w.headers {
		req.Header.Set(k, v)
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return true, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	// Drain body so the connection can be reused.
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return false, nil
	}

	retryable = resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500
	return retryable, fmt.Errorf("webhook returned status %d", resp.StatusCode)
}

func (w *WebhookNotifier) renderPayload(alert *domain.Alert) ([]byte, error) {
	if w.tmpl == nil {
		return json.Marshal(alert)
	}
	var buf bytes.Buffer
	if err := w.tmpl.Execute(&buf, alert); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (w *WebhookNotifier) Type() string {
	return "webhook"
}

func (w *WebhookNotifier) Validate(config map[string]interface{}) error {
	cfg, err := parseWebhookConfig(config)
	if err != nil {
		return err
	}
	if cfg.Payload != "" {
		if _, err := template.New("webhook").Parse(cfg.Payload); err != nil {
			return fmt.Errorf("invalid payload template: %w", err)
		}
	}
	return nil
}
