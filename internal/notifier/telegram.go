package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"health-monitor/internal/domain"
)

const (
	telegramAPIURL = "https://api.telegram.org/bot%s/sendMessage"
)

type TelegramNotifier struct {
	botToken string
	chatIDs  []string
	client   *http.Client
	log      zerolog.Logger
}

type telegramConfig struct {
	BotToken string   `json:"bot_token"`
	ChatIDs  []string `json:"chat_ids"`
	ProxyURL string   `json:"proxy_url"`
}

type telegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

type telegramResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description,omitempty"`
}

func NewTelegramNotifier(config map[string]interface{}, log zerolog.Logger) (*TelegramNotifier, error) {
	cfg, err := parseTelegramConfig(config)
	if err != nil {
		return nil, err
	}

	client, err := buildHTTPClient(cfg.ProxyURL, 10*time.Second)
	if err != nil {
		return nil, err
	}

	return &TelegramNotifier{
		botToken: cfg.BotToken,
		chatIDs:  cfg.ChatIDs,
		client:   client,
		log:      log,
	}, nil
}

func parseTelegramConfig(config map[string]interface{}) (*telegramConfig, error) {
	botToken, ok := config["bot_token"].(string)
	if !ok || botToken == "" {
		return nil, fmt.Errorf("bot_token is required")
	}

	chatIDsRaw, ok := config["chat_ids"]
	if !ok {
		return nil, fmt.Errorf("chat_ids is required")
	}

	var chatIDs []string
	switch v := chatIDsRaw.(type) {
	case []interface{}:
		for _, id := range v {
			if strID, ok := id.(string); ok {
				chatIDs = append(chatIDs, strID)
			}
		}
	case []string:
		chatIDs = v
	default:
		return nil, fmt.Errorf("chat_ids must be an array of strings")
	}

	if len(chatIDs) == 0 {
		return nil, fmt.Errorf("at least one chat_id is required")
	}

	var proxyURL string
	if raw, ok := config["proxy_url"]; ok {
		if proxyURL, ok = raw.(string); !ok {
			return nil, fmt.Errorf("proxy_url must be a string")
		}
		if proxyURL != "" {
			if _, err := parseProxyURL(proxyURL); err != nil {
				return nil, err
			}
		}
	}

	return &telegramConfig{
		BotToken: botToken,
		ChatIDs:  chatIDs,
		ProxyURL: proxyURL,
	}, nil
}

func (t *TelegramNotifier) Notify(ctx context.Context, alert *domain.Alert) error {
	message := t.formatMessage(alert)

	var lastErr error
	successCount := 0

	for _, chatID := range t.chatIDs {
		if err := t.sendMessage(ctx, chatID, message); err != nil {
			t.log.Error().
				Err(err).
				Str("chat_id", chatID).
				Str("alert_id", alert.ID).
				Msg("Failed to send Telegram notification")
			lastErr = err
		} else {
			successCount++
			t.log.Info().
				Str("chat_id", chatID).
				Str("alert_id", alert.ID).
				Str("alert_type", string(alert.Type)).
				Msg("Telegram notification sent")
		}
	}

	if successCount == 0 && lastErr != nil {
		return fmt.Errorf("failed to send to all chat IDs: %w", lastErr)
	}

	return nil
}

func (t *TelegramNotifier) sendMessage(ctx context.Context, chatID, text string) error {
	apiURL := fmt.Sprintf(telegramAPIURL, t.botToken)

	msg := telegramMessage{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "Markdown",
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var telegramResp telegramResponse
	if err := json.NewDecoder(resp.Body).Decode(&telegramResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !telegramResp.OK {
		return fmt.Errorf("telegram API error: %s", telegramResp.Description)
	}

	return nil
}

func (t *TelegramNotifier) formatMessage(alert *domain.Alert) string {
	var sb strings.Builder

	icon := t.getAlertIcon(alert.Type, alert.Severity)
	sb.WriteString(fmt.Sprintf("%s *%s*\n\n", icon, escapeMarkdown(alert.Message)))

	sb.WriteString(fmt.Sprintf("🎯 *Target:* %s\n", escapeMarkdown(alert.TargetName)))
	sb.WriteString(fmt.Sprintf("📊 *Type:* %s\n", escapeMarkdown(string(alert.Type))))
	sb.WriteString(fmt.Sprintf("🚨 *Severity:* %s\n", escapeMarkdown(string(alert.Severity))))

	if alert.Description != "" {
		sb.WriteString(fmt.Sprintf("📝 *Details:* %s\n", escapeMarkdown(alert.Description)))
	}

	if len(alert.Metadata) > 0 {
		sb.WriteString("\n📋 *Metadata:*\n")
		for key, value := range alert.Metadata {
			sb.WriteString(fmt.Sprintf("  • %s: %v\n", escapeMarkdown(key), value))
		}
	}

	sb.WriteString(fmt.Sprintf("\n🕐 *Time:* %s", alert.CreatedAt.Format("2006-01-02 15:04:05")))

	return sb.String()
}

func (t *TelegramNotifier) getAlertIcon(alertType domain.AlertType, severity domain.AlertSeverity) string {
	switch alertType {
	case domain.AlertTypeDown:
		return "🔴"
	case domain.AlertTypeUp:
		return "🟢"
	case domain.AlertTypeSlowResponse:
		return "🐌"
	case domain.AlertTypeSSLExpiring:
		return "🔐"
	case domain.AlertTypeConsecutiveFail:
		return "⚠️"
	default:
		if severity == domain.AlertSeverityCritical {
			return "🚨"
		}
		return "ℹ️"
	}
}

func escapeMarkdown(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)
	return replacer.Replace(text)
}

func (t *TelegramNotifier) Type() string {
	return "telegram"
}

func (t *TelegramNotifier) Validate(config map[string]interface{}) error {
	_, err := parseTelegramConfig(config)
	return err
}
