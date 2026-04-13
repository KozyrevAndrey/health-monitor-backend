package notifier

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"health-monitor/internal/domain"
)

// GmailNotifier sends alerts via Gmail API using a Service Account key file.
// The key file must be placed at the path specified by service_account_file.
// This path should point to a gitignored secrets/ directory.
type GmailNotifier struct {
	serviceAccountFile string
	from               string
	to                 []string
	impersonateUser    string // required for domain-wide delegation (Google Workspace)
	log                zerolog.Logger
}

type gmailConfig struct {
	ServiceAccountFile string   `json:"service_account_file"`
	From               string   `json:"from"`
	To                 []string `json:"to"`
	ImpersonateUser    string   `json:"impersonate_user"`
}

func NewGmailNotifier(config map[string]interface{}, log zerolog.Logger) (*GmailNotifier, error) {
	cfg, err := parseGmailConfig(config)
	if err != nil {
		return nil, err
	}

	return &GmailNotifier{
		serviceAccountFile: cfg.ServiceAccountFile,
		from:               cfg.From,
		to:                 cfg.To,
		impersonateUser:    cfg.ImpersonateUser,
		log:                log,
	}, nil
}

func parseGmailConfig(config map[string]interface{}) (*gmailConfig, error) {
	saFile, ok := config["service_account_file"].(string)
	if !ok || saFile == "" {
		return nil, fmt.Errorf("service_account_file is required")
	}

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
			if s, ok := addr.(string); ok {
				to = append(to, s)
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

	impersonateUser, _ := config["impersonate_user"].(string)

	return &gmailConfig{
		ServiceAccountFile: saFile,
		From:               from,
		To:                 to,
		ImpersonateUser:    impersonateUser,
	}, nil
}

func (g *GmailNotifier) Notify(ctx context.Context, alert *domain.Alert) error {
	svc, err := g.buildGmailService(ctx)
	if err != nil {
		return fmt.Errorf("failed to build Gmail service: %w", err)
	}

	subject := g.buildSubject(alert)
	htmlBody := g.buildHTMLBody(alert)
	plainBody := g.buildPlainBody(alert)

	raw, err := g.buildRawMessage(subject, htmlBody, plainBody)
	if err != nil {
		return fmt.Errorf("failed to build message: %w", err)
	}

	sendUser := "me"
	if g.impersonateUser != "" {
		sendUser = g.impersonateUser
	}

	_, err = svc.Users.Messages.Send(sendUser, &gmail.Message{Raw: raw}).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to send via Gmail API: %w", err)
	}

	g.log.Info().
		Str("alert_id", alert.ID).
		Str("alert_type", string(alert.Type)).
		Strs("recipients", g.to).
		Msg("Gmail notification sent")

	return nil
}

func (g *GmailNotifier) buildGmailService(ctx context.Context) (*gmail.Service, error) {
	saJSON, err := os.ReadFile(g.serviceAccountFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read service account file %s: %w", g.serviceAccountFile, err)
	}

	if g.impersonateUser != "" {
		// Domain-wide delegation: impersonate a real user in Google Workspace
		jwtCfg, err := google.JWTConfigFromJSON(saJSON, gmail.GmailSendScope)
		if err != nil {
			return nil, fmt.Errorf("failed to parse JWT config: %w", err)
		}
		jwtCfg.Subject = g.impersonateUser
		return gmail.NewService(ctx, option.WithHTTPClient(jwtCfg.Client(ctx)))
	}

	// No impersonation — service account acts as itself
	creds, err := google.CredentialsFromJSON(ctx, saJSON, gmail.GmailSendScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse service account credentials: %w", err)
	}

	return gmail.NewService(ctx, option.WithCredentials(creds))
}

func (g *GmailNotifier) buildRawMessage(subject, htmlBody, plainBody string) (string, error) {
	boundary := fmt.Sprintf("boundary_%d", time.Now().Unix())

	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("From: %s\r\n", g.from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(g.to, ", ")))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n", boundary))
	msg.WriteString("\r\n")

	msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
	msg.WriteString(plainBody)
	msg.WriteString("\r\n")

	msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
	msg.WriteString(htmlBody)
	msg.WriteString("\r\n")

	msg.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	return base64.URLEncoding.EncodeToString([]byte(msg.String())), nil
}

func (g *GmailNotifier) buildSubject(alert *domain.Alert) string {
	prefix := map[domain.AlertSeverity]string{
		domain.AlertSeverityCritical: "[CRITICAL]",
		domain.AlertSeverityWarning:  "[WARNING]",
		domain.AlertSeverityInfo:     "[INFO]",
	}
	p := prefix[alert.Severity]
	return fmt.Sprintf("%s %s - %s", p, alert.TargetName, strings.ToUpper(string(alert.Type)))
}

func (g *GmailNotifier) buildPlainBody(alert *domain.Alert) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s %s\n\n", gmailAlertIcon(alert.Type), alert.Message))
	sb.WriteString(fmt.Sprintf("Target:     %s\n", alert.TargetName))
	sb.WriteString(fmt.Sprintf("Alert Type: %s\n", alert.Type))
	sb.WriteString(fmt.Sprintf("Severity:   %s\n", alert.Severity))
	if alert.Description != "" {
		sb.WriteString(fmt.Sprintf("Details:    %s\n", alert.Description))
	}
	sb.WriteString(fmt.Sprintf("Time:       %s\n", alert.CreatedAt.Format("2006-01-02 15:04:05")))
	if len(alert.Metadata) > 0 {
		sb.WriteString("\nAdditional Information:\n")
		for k, v := range alert.Metadata {
			sb.WriteString(fmt.Sprintf("  %s: %v\n", k, v))
		}
	}
	sb.WriteString("\n---\nHealth Monitor\n")
	return sb.String()
}

func (g *GmailNotifier) buildHTMLBody(alert *domain.Alert) string {
	severityColors := map[domain.AlertSeverity]string{
		domain.AlertSeverityCritical: "#dc3545",
		domain.AlertSeverityWarning:  "#ffc107",
		domain.AlertSeverityInfo:     "#17a2b8",
	}
	color := severityColors[alert.Severity]
	if color == "" {
		color = "#6c757d"
	}

	type tmplData struct {
		Icon          string
		Message       string
		TargetName    string
		Type          string
		Severity      string
		Description   string
		CreatedAt     string
		Metadata      map[string]interface{}
		SeverityColor string
	}

	data := tmplData{
		Icon:          gmailAlertIcon(alert.Type),
		Message:       alert.Message,
		TargetName:    alert.TargetName,
		Type:          string(alert.Type),
		Severity:      string(alert.Severity),
		Description:   alert.Description,
		CreatedAt:     alert.CreatedAt.Format("2006-01-02 15:04:05"),
		Metadata:      alert.Metadata,
		SeverityColor: color,
	}

	const tmpl = `<!DOCTYPE html><html><head><meta charset="UTF-8">
<style>
body{font-family:Arial,sans-serif;line-height:1.6;color:#333;max-width:600px;margin:0 auto;padding:20px}
.header{background:{{.SeverityColor}};color:white;padding:20px;border-radius:5px 5px 0 0;text-align:center}
.header h1{margin:0;font-size:22px}
.content{border:1px solid #ddd;border-top:none;border-radius:0 0 5px 5px;padding:20px;background:#f9f9f9}
.row{margin:8px 0;padding:10px;background:white;border-radius:3px}
.label{font-weight:bold;color:#555;display:inline-block;min-width:110px}
.footer{margin-top:20px;text-align:center;font-size:12px;color:#888}
</style></head><body>
<div class="header"><h1>{{.Icon}} {{.Message}}</h1></div>
<div class="content">
<div class="row"><span class="label">Target:</span> {{.TargetName}}</div>
<div class="row"><span class="label">Alert Type:</span> {{.Type}}</div>
<div class="row"><span class="label">Severity:</span> {{.Severity}}</div>
{{if .Description}}<div class="row"><span class="label">Details:</span> {{.Description}}</div>{{end}}
<div class="row"><span class="label">Time:</span> {{.CreatedAt}}</div>
{{if .Metadata}}<div class="row"><strong>Additional Info:</strong>
{{range $k,$v := .Metadata}}<div style="padding:3px 0"><span class="label">{{$k}}:</span> {{$v}}</div>{{end}}
</div>{{end}}
</div>
<div class="footer">Health Monitor — Automated Alert System</div>
</body></html>`

	t := template.Must(template.New("gmail").Parse(tmpl))
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return g.buildPlainBody(alert)
	}
	return buf.String()
}

func (g *GmailNotifier) Type() string {
	return "gmail"
}

func (g *GmailNotifier) Validate(config map[string]interface{}) error {
	_, err := parseGmailConfig(config)
	return err
}

func gmailAlertIcon(t domain.AlertType) string {
	icons := map[domain.AlertType]string{
		domain.AlertTypeDown:            "🔴",
		domain.AlertTypeUp:              "🟢",
		domain.AlertTypeSlowResponse:    "🐌",
		domain.AlertTypeSSLExpiring:     "🔐",
		domain.AlertTypeConsecutiveFail: "⚠️",
	}
	if icon, ok := icons[t]; ok {
		return icon
	}
	return "ℹ️"
}
