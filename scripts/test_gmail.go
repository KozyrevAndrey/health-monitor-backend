//go:build ignore

// Test script for Gmail API sending via Service Account.
//
// Usage:
//   go run scripts/test_gmail.go
//   go run scripts/test_gmail.go -sa secrets/google-service-account.json -to andrey0kozyrev0@gmail.com
//
// Flags:
//   -sa   path to service account JSON file (default: secrets/google-service-account.json)
//   -from sender address (default: andrey0kozyrev0@gmail.com)
//   -to   recipient address (default: andrey0kozyrev0@gmail.com)
//   -imp  impersonate user for domain-wide delegation (optional)

package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func main() {
	saFile  := flag.String("sa",   "secrets/google-service-account.json", "Service account JSON file path")
	from    := flag.String("from", "andrey0kozyrev0@gmail.com",           "From address")
	to      := flag.String("to",   "andrey0kozyrev0@gmail.com",           "Recipient address")
	imp     := flag.String("imp",  "",                                     "Impersonate user (domain-wide delegation)")
	flag.Parse()

	fmt.Printf("📧 Gmail API Test\n")
	fmt.Printf("   SA file : %s\n", *saFile)
	fmt.Printf("   From    : %s\n", *from)
	fmt.Printf("   To      : %s\n", *to)
	if *imp != "" {
		fmt.Printf("   Imp     : %s\n", *imp)
	}
	fmt.Println()

	saJSON, err := os.ReadFile(*saFile)
	if err != nil {
		log.Fatalf("❌ Failed to read service account file: %v\n\n"+
			"Make sure the file exists:\n  mkdir -p secrets\n  cp your-sa-file.json %s\n", err, *saFile)
	}
	fmt.Printf("✅ Service account file loaded (%d bytes)\n", len(saJSON))

	ctx := context.Background()
	svc, err := buildService(ctx, saJSON, *imp)
	if err != nil {
		log.Fatalf("❌ Failed to create Gmail service: %v", err)
	}
	fmt.Println("✅ Gmail service created")

	msg := buildMessage(*from, *to)
	sendUser := "me"
	if *imp != "" {
		sendUser = *imp
	}

	fmt.Printf("📤 Sending test email to %s...\n", *to)
	_, err = svc.Users.Messages.Send(sendUser, &gmail.Message{Raw: msg}).Context(ctx).Do()
	if err != nil {
		log.Fatalf("❌ Failed to send email: %v", err)
	}

	fmt.Println("✅ Email sent successfully! Check your inbox.")
}

func buildService(ctx context.Context, saJSON []byte, impersonateUser string) (*gmail.Service, error) {
	if impersonateUser != "" {
		jwtCfg, err := google.JWTConfigFromJSON(saJSON, gmail.GmailSendScope)
		if err != nil {
			return nil, fmt.Errorf("JWT config: %w", err)
		}
		jwtCfg.Subject = impersonateUser
		return gmail.NewService(ctx, option.WithHTTPClient(jwtCfg.Client(ctx)))
	}

	creds, err := google.CredentialsFromJSON(ctx, saJSON, gmail.GmailSendScope)
	if err != nil {
		return nil, fmt.Errorf("credentials: %w", err)
	}
	return gmail.NewService(ctx, option.WithCredentials(creds))
}

func buildMessage(from, to string) string {
	now := time.Now().Format("2006-01-02 15:04:05")
	boundary := fmt.Sprintf("boundary_%d", time.Now().Unix())

	plain := fmt.Sprintf("Health Monitor — Test Email\n\nThis is a test notification sent at %s.\nIf you received this, Gmail API is working correctly!\n", now)

	html := fmt.Sprintf(`<!DOCTYPE html><html><body style="font-family:Arial,sans-serif;max-width:500px;margin:40px auto;padding:20px">
<div style="background:linear-gradient(135deg,#667eea,#764ba2);color:white;padding:24px;border-radius:8px 8px 0 0;text-align:center">
<h1 style="margin:0;font-size:22px">✅ Health Monitor</h1>
<p style="margin:8px 0 0;opacity:.9">Gmail API Test</p>
</div>
<div style="border:1px solid #ddd;border-top:none;border-radius:0 0 8px 8px;padding:24px;background:#f9f9f9">
<p>Gmail API is configured correctly and working!</p>
<table style="width:100%%;border-collapse:collapse">
<tr><td style="padding:8px;background:white;border-radius:4px;font-weight:bold;width:100px">From:</td>
    <td style="padding:8px;background:white;border-radius:4px">%s</td></tr>
<tr><td style="padding:8px;font-weight:bold">Time:</td>
    <td style="padding:8px">%s</td></tr>
</table>
</div>
<p style="text-align:center;font-size:12px;color:#888;margin-top:16px">Health Monitor — Automated Alert System</p>
</body></html>`, from, now)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("From: %s\r\n", from))
	sb.WriteString(fmt.Sprintf("To: %s\r\n", to))
	sb.WriteString("Subject: [Health Monitor] Test Email — Gmail API OK\r\n")
	sb.WriteString("MIME-Version: 1.0\r\n")
	sb.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n\r\n", boundary))
	sb.WriteString(fmt.Sprintf("--%s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s\r\n", boundary, plain))
	sb.WriteString(fmt.Sprintf("--%s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s\r\n", boundary, html))
	sb.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	return base64.URLEncoding.EncodeToString([]byte(sb.String()))
}
