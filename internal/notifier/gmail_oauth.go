package notifier

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"health-monitor/internal/domain"
)

// GmailOAuthNotifier sends alerts via Gmail API using OAuth2 (personal Gmail account).
// Requires two files in secrets/:
//   - credentials file: OAuth2 client credentials (downloaded from Google Cloud Console)
//   - token file: generated once by running scripts/gmail_auth.go
type GmailOAuthNotifier struct {
	credentialsFile string
	tokenFile       string
	from            string
	to              []string
	log             zerolog.Logger
}

type gmailOAuthConfig struct {
	CredentialsFile string   `json:"credentials_file"`
	TokenFile       string   `json:"token_file"`
	From            string   `json:"from"`
	To              []string `json:"to"`
}

func NewGmailOAuthNotifier(config map[string]interface{}, log zerolog.Logger) (*GmailOAuthNotifier, error) {
	cfg, err := parseGmailOAuthConfig(config)
	if err != nil {
		return nil, err
	}
	return &GmailOAuthNotifier{
		credentialsFile: cfg.CredentialsFile,
		tokenFile:       cfg.TokenFile,
		from:            cfg.From,
		to:              cfg.To,
		log:             log,
	}, nil
}

func parseGmailOAuthConfig(config map[string]interface{}) (*gmailOAuthConfig, error) {
	credFile, _ := config["credentials_file"].(string)
	if credFile == "" {
		credFile = "secrets/gmail-token.json"
	}

	tokenFile, _ := config["token_file"].(string)
	if tokenFile == "" {
		tokenFile = "secrets/gmail-oauth-token.json"
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

	return &gmailOAuthConfig{
		CredentialsFile: credFile,
		TokenFile:       tokenFile,
		From:            from,
		To:              to,
	}, nil
}

func (g *GmailOAuthNotifier) Notify(ctx context.Context, alert *domain.Alert) error {
	svc, err := g.buildGmailService(ctx)
	if err != nil {
		return fmt.Errorf("failed to build Gmail service: %w", err)
	}

	subject  := g.buildSubject(alert)
	htmlBody := g.buildHTMLBody(alert)
	plainBody := g.buildPlainBody(alert)

	raw, err := g.buildRawMessage(subject, htmlBody, plainBody)
	if err != nil {
		return fmt.Errorf("failed to build message: %w", err)
	}

	_, err = svc.Users.Messages.Send("me", &gmail.Message{Raw: raw}).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to send via Gmail API: %w", err)
	}

	g.log.Info().
		Str("alert_id", alert.ID).
		Str("alert_type", string(alert.Type)).
		Strs("recipients", g.to).
		Msg("Gmail OAuth notification sent")

	return nil
}

func (g *GmailOAuthNotifier) buildGmailService(ctx context.Context) (*gmail.Service, error) {
	credJSON, err := os.ReadFile(g.credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("read credentials file %s: %w", g.credentialsFile, err)
	}

	oauthCfg, err := google.ConfigFromJSON(credJSON, gmail.GmailSendScope)
	if err != nil {
		return nil, fmt.Errorf("parse credentials: %w", err)
	}

	token, err := g.loadToken()
	if err != nil {
		return nil, fmt.Errorf("load token from %s: %w\n"+
			"Run: go run scripts/gmail_auth.go", g.tokenFile, err)
	}

	// TokenSource automatically refreshes the access_token using the refresh_token
	tokenSource := oauthCfg.TokenSource(ctx, token)

	// Persist refreshed token back to file if it changed
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("refresh token: %w", err)
	}
	if newToken.AccessToken != token.AccessToken {
		if saveErr := g.saveToken(newToken); saveErr != nil {
			g.log.Warn().Err(saveErr).Msg("Failed to persist refreshed token")
		}
	}

	return gmail.NewService(ctx, option.WithTokenSource(tokenSource))
}

func (g *GmailOAuthNotifier) loadToken() (*oauth2.Token, error) {
	f, err := os.Open(g.tokenFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var token oauth2.Token
	if err := json.NewDecoder(f).Decode(&token); err != nil {
		return nil, err
	}
	return &token, nil
}

func (g *GmailOAuthNotifier) saveToken(token *oauth2.Token) error {
	f, err := os.OpenFile(g.tokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}

func (g *GmailOAuthNotifier) buildRawMessage(subject, htmlBody, plainBody string) (string, error) {
	boundary := fmt.Sprintf("boundary_%d", time.Now().Unix())
	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("From: %s\r\n", g.from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(g.to, ", ")))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n\r\n", boundary))
	msg.WriteString(fmt.Sprintf("--%s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s\r\n", boundary, plainBody))
	msg.WriteString(fmt.Sprintf("--%s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s\r\n", boundary, htmlBody))
	msg.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	return base64.URLEncoding.EncodeToString([]byte(msg.String())), nil
}

func (g *GmailOAuthNotifier) buildSubject(alert *domain.Alert) string {
	prefix := map[domain.AlertSeverity]string{
		domain.AlertSeverityCritical: "[CRITICAL]",
		domain.AlertSeverityWarning:  "[WARNING]",
		domain.AlertSeverityInfo:     "[INFO]",
	}
	return fmt.Sprintf("%s %s - %s", prefix[alert.Severity], alert.TargetName, strings.ToUpper(string(alert.Type)))
}

func (g *GmailOAuthNotifier) buildPlainBody(alert *domain.Alert) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s %s\n\n", oauthAlertIcon(alert.Type), alert.Message))
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

func (g *GmailOAuthNotifier) buildHTMLBody(alert *domain.Alert) string {
	colors := map[domain.AlertSeverity]string{
		domain.AlertSeverityCritical: "#dc3545",
		domain.AlertSeverityWarning:  "#ffc107",
		domain.AlertSeverityInfo:     "#17a2b8",
	}
	color := colors[alert.Severity]
	if color == "" {
		color = "#6c757d"
	}

	type data struct {
		Icon, Message, TargetName, Type, Severity, Description, CreatedAt, SeverityColor string
		Metadata map[string]interface{}
	}

	d := data{
		Icon:          oauthAlertIcon(alert.Type),
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

	t := template.Must(template.New("gmail_oauth").Parse(tmpl))
	var buf bytes.Buffer
	if err := t.Execute(&buf, d); err != nil {
		return g.buildPlainBody(alert)
	}
	return buf.String()
}

func (g *GmailOAuthNotifier) Type() string {
	return "gmail_oauth"
}

func (g *GmailOAuthNotifier) Validate(config map[string]interface{}) error {
	_, err := parseGmailOAuthConfig(config)
	return err
}

func oauthAlertIcon(t domain.AlertType) string {
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
