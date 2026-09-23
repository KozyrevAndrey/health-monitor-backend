package notifier

import (
	"fmt"

	"github.com/rs/zerolog"

	"health-monitor/internal/domain"
)

// New builds a notifier from its stored config. Both the startup loader and the
// API reload path go through here, so adding a provider is one case in one
// switch instead of two that drift apart.
func New(cfg *domain.NotifierConfig, log zerolog.Logger) (domain.Notifier, error) {
	switch cfg.Type {
	case "telegram":
		return NewTelegramNotifier(cfg.Config, log)
	case "email":
		return NewEmailNotifier(cfg.Config, log)
	case "gmail":
		return NewGmailNotifier(cfg.Config, log)
	case "gmail_oauth":
		return NewGmailOAuthNotifier(cfg.Config, log)
	case "webhook":
		return NewWebhookNotifier(cfg.Config, log)
	case domain.NotifierTypeSMSPilot:
		return NewSMSPilotNotifier(cfg.Config, log)
	default:
		return nil, fmt.Errorf("unknown notifier type: %s", cfg.Type)
	}
}
