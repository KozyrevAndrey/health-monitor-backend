package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/rs/zerolog"
	"health-monitor/internal/domain"
)

const (
	smspilotDefaultAPIURL   = "https://smspilot.ru/api.php"
	smspilotDefaultKeyEnv   = "SMSPILOT_API_KEY"
	smspilotAPIURLEnv       = "SMSPILOT_API_URL"
	smspilotDefaultVoice    = "GOLOS"
	smspilotTimeout         = 20 * time.Second
	smspilotDefaultTemplate = "Внимание! {{.TargetName}} недоступен. Проверьте мониторинг."
)

// smspilotPhone is the only format SMSPILOT accepts for Russian mobile numbers.
var smspilotPhone = regexp.MustCompile(`^79[0-9]{9}$`)

// SMSPilotNotifier places voice calls through SMSPILOT. Unlike chat notifiers it
// is policy-routed (see domain.IsPolicyRouted): the alert manager decides when a
// call is allowed, this type only performs it.
type SMSPilotNotifier struct {
	apiURL string
	apiKey string
	phones []string
	voice  string
	ttl    int
	tmpl   *template.Template
	client *http.Client
	log    zerolog.Logger
}

type smspilotConfig struct {
	Phones   []string
	Voice    string
	TTL      int
	Template string
	KeyEnv   string
}

// smspilotResponse covers both shapes SMSPILOT returns: an error object, or a
// list with one item per number.
type smspilotResponse struct {
	Error *smspilotError `json:"error"`
	Send  []smspilotItem `json:"send"`
}

type smspilotError struct {
	Code          interface{} `json:"code"`
	Description   string      `json:"description"`
	DescriptionRu string      `json:"description_ru"`
}

type smspilotItem struct {
	ServerID interface{} `json:"server_id"`
	Status   interface{} `json:"status"`
	Phone    string      `json:"phone"`
	Price    interface{} `json:"price"`
	Error    string      `json:"error"`
	ErrorRu  string      `json:"error_ru"`
}

// NewSMSPilotNotifier builds the notifier against the real SMSPILOT endpoint.
func NewSMSPilotNotifier(config map[string]interface{}, log zerolog.Logger) (*SMSPilotNotifier, error) {
	return newSMSPilotNotifier(smspilotAPIURL(), config, log)
}

// newSMSPilotNotifier is the injectable constructor. The endpoint is never read
// from the notifier config: config keys are merged from API requests, so a
// dashboard user could otherwise point the request — API key included — at a
// host of their choosing. Tests pass an httptest URL here, and a server-side
// SMSPILOT_API_URL override exists for manual end-to-end checks.
func newSMSPilotNotifier(apiURL string, config map[string]interface{}, log zerolog.Logger) (*SMSPilotNotifier, error) {
	cfg, err := parseSMSPilotConfig(config)
	if err != nil {
		return nil, err
	}

	// The delivery policy lives in the same config blob and is parsed by the
	// alert manager. Rejecting a bad one here is what keeps the API from saving
	// a channel that looks enabled in the dashboard but never registers.
	if _, err := domain.ParseDeliveryPolicy(config); err != nil {
		return nil, err
	}

	apiKey := os.Getenv(cfg.KeyEnv)
	if apiKey == "" {
		return nil, fmt.Errorf("environment variable %s is not set on the server", cfg.KeyEnv)
	}

	tmpl, err := template.New("smspilot").Parse(cfg.Template)
	if err != nil {
		return nil, fmt.Errorf("invalid message_template: %w", err)
	}

	client, err := buildHTTPClient("", smspilotTimeout)
	if err != nil {
		return nil, err
	}

	return &SMSPilotNotifier{
		apiURL: apiURL,
		apiKey: apiKey,
		phones: cfg.Phones,
		voice:  cfg.Voice,
		ttl:    cfg.TTL,
		tmpl:   tmpl,
		client: client,
		log:    log,
	}, nil
}

func smspilotAPIURL() string {
	if override := os.Getenv(smspilotAPIURLEnv); override != "" {
		return override
	}
	return smspilotDefaultAPIURL
}

func parseSMSPilotConfig(config map[string]interface{}) (*smspilotConfig, error) {
	cfg := &smspilotConfig{
		Voice:    smspilotDefaultVoice,
		Template: smspilotDefaultTemplate,
		KeyEnv:   smspilotDefaultKeyEnv,
	}

	phones, err := smspilotPhones(config["phones"])
	if err != nil {
		return nil, err
	}
	cfg.Phones = phones

	if raw, ok := config["voice"]; ok {
		voice, ok := raw.(string)
		if !ok {
			return nil, fmt.Errorf("voice must be a string")
		}
		if voice != "" {
			voice = strings.ToUpper(voice)
			if voice != "GOLOS" && voice != "GOLOSM" {
				return nil, fmt.Errorf("voice must be GOLOS or GOLOSM")
			}
			cfg.Voice = voice
		}
	}

	if raw, ok := config["message_template"]; ok {
		tmpl, ok := raw.(string)
		if !ok {
			return nil, fmt.Errorf("message_template must be a string")
		}
		if strings.TrimSpace(tmpl) != "" {
			cfg.Template = tmpl
		}
	}

	if raw, ok := config["ttl"]; ok {
		ttl, err := toInt(raw)
		if err != nil {
			return nil, fmt.Errorf("ttl must be a number")
		}
		if ttl < 0 {
			return nil, fmt.Errorf("ttl must be >= 0")
		}
		cfg.TTL = ttl
	}

	if raw, ok := config["api_key_env"]; ok {
		keyEnv, ok := raw.(string)
		if !ok {
			return nil, fmt.Errorf("api_key_env must be a string")
		}
		if keyEnv != "" {
			cfg.KeyEnv = keyEnv
		}
	}

	return cfg, nil
}

func smspilotPhones(raw interface{}) ([]string, error) {
	if raw == nil {
		return nil, fmt.Errorf("phones is required")
	}

	var values []string
	switch v := raw.(type) {
	case []string:
		values = v
	case []interface{}:
		for _, item := range v {
			phone, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("phones must be an array of strings")
			}
			values = append(values, phone)
		}
	default:
		return nil, fmt.Errorf("phones must be an array of strings")
	}

	phones := make([]string, 0, len(values))
	for _, value := range values {
		phone := strings.TrimPrefix(strings.ReplaceAll(strings.TrimSpace(value), " ", ""), "+")
		if !smspilotPhone.MatchString(phone) {
			return nil, fmt.Errorf("phone %q must look like 79XXXXXXXXX", value)
		}
		phones = append(phones, phone)
	}

	if len(phones) == 0 {
		return nil, fmt.Errorf("at least one phone is required")
	}

	return phones, nil
}

func (s *SMSPilotNotifier) Type() string {
	return domain.NotifierTypeSMSPilot
}

func (s *SMSPilotNotifier) Validate(config map[string]interface{}) error {
	cfg, err := parseSMSPilotConfig(config)
	if err != nil {
		return err
	}

	if _, err := template.New("smspilot").Parse(cfg.Template); err != nil {
		return fmt.Errorf("invalid message_template: %w", err)
	}

	if os.Getenv(cfg.KeyEnv) == "" {
		return fmt.Errorf("environment variable %s is not set on the server", cfg.KeyEnv)
	}

	if _, err := domain.ParseDeliveryPolicy(config); err != nil {
		return err
	}

	return nil
}

// Notify places the call and discards provider identifiers.
func (s *SMSPilotNotifier) Notify(ctx context.Context, alert *domain.Alert) error {
	_, err := s.NotifyWithReceipt(ctx, alert)
	return err
}

// NotifyWithReceipt places one call per configured phone and returns the
// provider identifiers of the calls that were accepted.
//
// Errors are split in two: an explicit rejection means nothing was dialled,
// while domain.ErrUnknownResult means the call may well have gone out and must
// never be retried automatically.
func (s *SMSPilotNotifier) NotifyWithReceipt(ctx context.Context, alert *domain.Alert) (domain.Receipt, error) {
	text, err := s.renderMessage(alert)
	if err != nil {
		return domain.Receipt{}, fmt.Errorf("failed to render message: %w", err)
	}

	var (
		serverIDs []string
		statuses  []string
		prices    []string
		firstErr  error
	)

	for _, phone := range s.phones {
		item, err := s.call(ctx, phone, text)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}

		serverIDs = append(serverIDs, smspilotString(item.ServerID))
		statuses = append(statuses, smspilotString(item.Status))
		prices = append(prices, smspilotString(item.Price))

		s.log.Info().
			Str("alert_id", alert.ID).
			Str("server_id", smspilotString(item.ServerID)).
			Str("status", smspilotString(item.Status)).
			Msg("SMSPILOT voice call accepted")
	}

	receipt := domain.Receipt{
		ExternalID:     strings.Join(serverIDs, ","),
		ExternalStatus: strings.Join(statuses, ","),
		Cost:           strings.Join(prices, ","),
	}

	return receipt, firstErr
}

// call places a single call. The request body carries the API key, so it is
// never logged.
func (s *SMSPilotNotifier) call(ctx context.Context, phone, text string) (*smspilotItem, error) {
	form := url.Values{}
	form.Set("apikey", s.apiKey)
	form.Set("format", "json")
	form.Set("to", phone)
	form.Set("send", text)
	form.Set("from", s.voice)
	if s.ttl > 0 {
		form.Set("ttl", strconv.Itoa(s.ttl))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: request to SMSPILOT failed: %v", domain.ErrUnknownResult, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to read SMSPILOT response: %v", domain.ErrUnknownResult, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: SMSPILOT returned HTTP %d", domain.ErrUnknownResult, resp.StatusCode)
	}

	var parsed smspilotResponse
	if err := json.Unmarshal(bytes.TrimSpace(body), &parsed); err != nil {
		return nil, fmt.Errorf("%w: unparseable SMSPILOT response", domain.ErrUnknownResult)
	}

	if parsed.Error != nil {
		description := parsed.Error.DescriptionRu
		if description == "" {
			description = parsed.Error.Description
		}
		return nil, fmt.Errorf("SMSPILOT rejected the call (%s): %s", smspilotString(parsed.Error.Code), description)
	}

	if len(parsed.Send) != 1 {
		return nil, fmt.Errorf("%w: SMSPILOT returned no unambiguous result", domain.ErrUnknownResult)
	}

	item := parsed.Send[0]

	if smspilotString(item.ServerID) == "" {
		return nil, fmt.Errorf("%w: SMSPILOT response has no server_id", domain.ErrUnknownResult)
	}

	status, err := strconv.Atoi(smspilotString(item.Status))
	if err != nil {
		return nil, fmt.Errorf("%w: SMSPILOT response has no numeric status", domain.ErrUnknownResult)
	}

	if status < 0 {
		reason := item.ErrorRu
		if reason == "" {
			reason = item.Error
		}
		return nil, fmt.Errorf("SMSPILOT rejected the call (status %d): %s", status, reason)
	}

	return &item, nil
}

func (s *SMSPilotNotifier) renderMessage(alert *domain.Alert) (string, error) {
	var buf bytes.Buffer
	if err := s.tmpl.Execute(&buf, alert); err != nil {
		return "", err
	}

	text := strings.TrimSpace(buf.String())
	if text == "" {
		return "", fmt.Errorf("message must not be empty")
	}

	return text, nil
}

// smspilotString normalises fields SMSPILOT sends as either a string or a number.
func smspilotString(raw interface{}) string {
	switch v := raw.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}
