package notifier

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"health-monitor/internal/domain"
)

type EmailNotifier struct {
	smtpHost     string
	smtpPort     int
	smtpUser     string
	smtpPassword string
	from         string
	to           []string
	useTLS       bool
	log          zerolog.Logger
}

type emailConfig struct {
	SMTPHost     string   `json:"smtp_host"`
	SMTPPort     int      `json:"smtp_port"`
	SMTPUser     string   `json:"smtp_user"`
	SMTPPassword string   `json:"smtp_password"`
	From         string   `json:"from"`
	To           []string `json:"to"`
	UseTLS       bool     `json:"use_tls"`
}

type emailData struct {
	Alert         *domain.Alert
	TargetName    string
	Type          string
	Severity      string
	Message       string
	Description   string
	CreatedAt     string
	Metadata      map[string]interface{}
	SeverityColor string
	Icon          string
}

func NewEmailNotifier(config map[string]interface{}, log zerolog.Logger) (*EmailNotifier, error) {
	cfg, err := parseEmailConfig(config)
	if err != nil {
		return nil, err
	}

	return &EmailNotifier{
		smtpHost:     cfg.SMTPHost,
		smtpPort:     cfg.SMTPPort,
		smtpUser:     cfg.SMTPUser,
		smtpPassword: cfg.SMTPPassword,
		from:         cfg.From,
		to:           cfg.To,
		useTLS:       cfg.UseTLS,
		log:          log,
	}, nil
}

func parseEmailConfig(config map[string]interface{}) (*emailConfig, error) {
	smtpHost, ok := config["smtp_host"].(string)
	if !ok || smtpHost == "" {
		return nil, fmt.Errorf("smtp_host is required")
	}

	var smtpPort int
	if port, ok := config["smtp_port"].(int); ok {
		smtpPort = port
	} else if port, ok := config["smtp_port"].(float64); ok {
		smtpPort = int(port)
	} else {
		smtpPort = 587
	}

	smtpUser, _ := config["smtp_user"].(string)
	smtpPassword, _ := config["smtp_password"].(string)

	from, ok := config["from"].(string)
	if !ok || from == "" {
		return nil, fmt.Errorf("from address is required")
	}

	toRaw, ok := config["to"]
	if !ok {
		return nil, fmt.Errorf("to addresses are required")
	}

	var to []string
	switch v := toRaw.(type) {
	case []interface{}:
		for _, addr := range v {
			if strAddr, ok := addr.(string); ok {
				to = append(to, strAddr)
			}
		}
	case []string:
		to = v
	default:
		return nil, fmt.Errorf("to must be an array of strings")
	}

	if len(to) == 0 {
		return nil, fmt.Errorf("at least one to address is required")
	}

	useTLS := true
	if val, ok := config["use_tls"].(bool); ok {
		useTLS = val
	}

	return &emailConfig{
		SMTPHost:     smtpHost,
		SMTPPort:     smtpPort,
		SMTPUser:     smtpUser,
		SMTPPassword: smtpPassword,
		From:         from,
		To:           to,
		UseTLS:       useTLS,
	}, nil
}

func (e *EmailNotifier) Notify(ctx context.Context, alert *domain.Alert) error {
	subject := e.buildSubject(alert)
	htmlBody := e.buildHTMLBody(alert)
	plainBody := e.buildPlainTextBody(alert)

	msg := e.buildEmailMessage(subject, htmlBody, plainBody)

	addr := fmt.Sprintf("%s:%d", e.smtpHost, e.smtpPort)
	auth := smtp.PlainAuth("", e.smtpUser, e.smtpPassword, e.smtpHost)

	err := smtp.SendMail(addr, auth, e.from, e.to, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	e.log.Info().
		Str("alert_id", alert.ID).
		Str("alert_type", string(alert.Type)).
		Strs("recipients", e.to).
		Msg("Email notification sent")

	return nil
}

func (e *EmailNotifier) buildSubject(alert *domain.Alert) string {
	severityPrefix := ""
	switch alert.Severity {
	case domain.AlertSeverityCritical:
		severityPrefix = "[CRITICAL]"
	case domain.AlertSeverityWarning:
		severityPrefix = "[WARNING]"
	case domain.AlertSeverityInfo:
		severityPrefix = "[INFO]"
	}

	return fmt.Sprintf("%s %s - %s", severityPrefix, alert.TargetName, strings.ToUpper(string(alert.Type)))
}

func (e *EmailNotifier) buildHTMLBody(alert *domain.Alert) string {
	data := e.prepareEmailData(alert)

	tmpl := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .header {
            background: {{.SeverityColor}};
            color: white;
            padding: 20px;
            border-radius: 5px 5px 0 0;
            text-align: center;
        }
        .header h1 {
            margin: 0;
            font-size: 24px;
        }
        .content {
            border: 1px solid #ddd;
            border-top: none;
            border-radius: 0 0 5px 5px;
            padding: 20px;
            background: #f9f9f9;
        }
        .info-row {
            margin: 10px 0;
            padding: 10px;
            background: white;
            border-radius: 3px;
        }
        .label {
            font-weight: bold;
            color: #555;
            display: inline-block;
            min-width: 120px;
        }
        .metadata {
            background: white;
            padding: 15px;
            border-radius: 3px;
            margin-top: 15px;
        }
        .metadata-item {
            padding: 5px 0;
            border-bottom: 1px solid #eee;
        }
        .metadata-item:last-child {
            border-bottom: none;
        }
        .footer {
            margin-top: 20px;
            padding-top: 20px;
            border-top: 1px solid #ddd;
            text-align: center;
            font-size: 12px;
            color: #888;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>{{.Icon}} {{.Message}}</h1>
    </div>
    <div class="content">
        <div class="info-row">
            <span class="label">Target:</span>
            <span>{{.TargetName}}</span>
        </div>
        <div class="info-row">
            <span class="label">Alert Type:</span>
            <span>{{.Type}}</span>
        </div>
        <div class="info-row">
            <span class="label">Severity:</span>
            <span>{{.Severity}}</span>
        </div>
        {{if .Description}}
        <div class="info-row">
            <span class="label">Details:</span>
            <span>{{.Description}}</span>
        </div>
        {{end}}
        <div class="info-row">
            <span class="label">Time:</span>
            <span>{{.CreatedAt}}</span>
        </div>
        {{if .Metadata}}
        <div class="metadata">
            <strong>Additional Information:</strong>
            {{range $key, $value := .Metadata}}
            <div class="metadata-item">
                <span class="label">{{$key}}:</span>
                <span>{{$value}}</span>
            </div>
            {{end}}
        </div>
        {{end}}
    </div>
    <div class="footer">
        <p>Health Monitor - Automated Alert System</p>
    </div>
</body>
</html>`

	t := template.Must(template.New("email").Parse(tmpl))
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		e.log.Error().Err(err).Msg("Failed to render HTML email template")
		return e.buildPlainTextBody(alert)
	}

	return buf.String()
}

func (e *EmailNotifier) buildPlainTextBody(alert *domain.Alert) string {
	var sb strings.Builder

	icon := e.getAlertIcon(alert.Type)
	sb.WriteString(fmt.Sprintf("%s %s\n\n", icon, alert.Message))
	sb.WriteString(fmt.Sprintf("Target: %s\n", alert.TargetName))
	sb.WriteString(fmt.Sprintf("Alert Type: %s\n", alert.Type))
	sb.WriteString(fmt.Sprintf("Severity: %s\n", alert.Severity))

	if alert.Description != "" {
		sb.WriteString(fmt.Sprintf("Details: %s\n", alert.Description))
	}

	sb.WriteString(fmt.Sprintf("Time: %s\n", alert.CreatedAt.Format("2006-01-02 15:04:05")))

	if len(alert.Metadata) > 0 {
		sb.WriteString("\nAdditional Information:\n")
		for key, value := range alert.Metadata {
			sb.WriteString(fmt.Sprintf("  %s: %v\n", key, value))
		}
	}

	sb.WriteString("\n---\nHealth Monitor - Automated Alert System\n")

	return sb.String()
}

func (e *EmailNotifier) prepareEmailData(alert *domain.Alert) *emailData {
	severityColor := e.getSeverityColor(alert.Severity)
	icon := e.getAlertIcon(alert.Type)

	return &emailData{
		Alert:         alert,
		TargetName:    alert.TargetName,
		Type:          string(alert.Type),
		Severity:      string(alert.Severity),
		Message:       alert.Message,
		Description:   alert.Description,
		CreatedAt:     alert.CreatedAt.Format("2006-01-02 15:04:05"),
		Metadata:      alert.Metadata,
		SeverityColor: severityColor,
		Icon:          icon,
	}
}

func (e *EmailNotifier) getSeverityColor(severity domain.AlertSeverity) string {
	switch severity {
	case domain.AlertSeverityCritical:
		return "#dc3545"
	case domain.AlertSeverityWarning:
		return "#ffc107"
	case domain.AlertSeverityInfo:
		return "#17a2b8"
	default:
		return "#6c757d"
	}
}

func (e *EmailNotifier) getAlertIcon(alertType domain.AlertType) string {
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
		return "ℹ️"
	}
}

func (e *EmailNotifier) buildEmailMessage(subject, htmlBody, plainBody string) string {
	boundary := fmt.Sprintf("boundary_%d", time.Now().Unix())

	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("From: %s\r\n", e.from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(e.to, ", ")))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n", boundary))
	msg.WriteString("\r\n")

	msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(plainBody)
	msg.WriteString("\r\n")

	msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)
	msg.WriteString("\r\n")

	msg.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	return msg.String()
}

func (e *EmailNotifier) Type() string {
	return "email"
}

func (e *EmailNotifier) Validate(config map[string]interface{}) error {
	_, err := parseEmailConfig(config)
	return err
}
