//go:build ignore

// Test script for Gmail OAuth2 notifier.
// Requires secrets/gmail-token.json (credentials) and secrets/gmail-oauth-token.json (token).
//
// Usage:
//   go run scripts/test_gmail_oauth.go
//   go run scripts/test_gmail_oauth.go -to other@gmail.com

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func main() {
	credFile  := flag.String("credentials", "secrets/gmail-token.json",       "OAuth2 client credentials JSON")
	tokenFile := flag.String("token",       "secrets/gmail-oauth-token.json", "Token file (from gmail_auth.go)")
	from      := flag.String("from",        "andrey0kozyrev0@gmail.com",       "From address")
	to        := flag.String("to",          "andrey0kozyrev0@gmail.com",       "Recipient address")
	flag.Parse()

	fmt.Println("📧 Gmail OAuth2 Test")
	fmt.Printf("   Credentials : %s\n", *credFile)
	fmt.Printf("   Token       : %s\n", *tokenFile)
	fmt.Printf("   From        : %s\n", *from)
	fmt.Printf("   To          : %s\n\n", *to)

	// Load credentials
	credJSON, err := os.ReadFile(*credFile)
	if err != nil {
		log.Fatalf("❌ Cannot read credentials: %v\n\nMake sure %s exists.", err, *credFile)
	}
	fmt.Println("✅ Credentials loaded")

	oauthCfg, err := google.ConfigFromJSON(credJSON, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("❌ Invalid credentials JSON: %v", err)
	}

	// Load token
	tokenJSON, err := os.ReadFile(*tokenFile)
	if err != nil {
		log.Fatalf("❌ Cannot read token: %v\n\nRun authorization first:\n  go run scripts/gmail_auth.go", err)
	}

	var token oauth2.Token
	if err := json.Unmarshal(tokenJSON, &token); err != nil {
		log.Fatalf("❌ Invalid token JSON: %v", err)
	}
	fmt.Printf("✅ Token loaded (expires: %s)\n", token.Expiry.Format("2006-01-02 15:04:05"))

	ctx := context.Background()
	tokenSource := oauthCfg.TokenSource(ctx, &token)

	svc, err := gmail.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		log.Fatalf("❌ Failed to create Gmail service: %v", err)
	}
	fmt.Println("✅ Gmail service created")

	// Build and send test message
	raw := buildTestMessage(*from, *to)
	fmt.Printf("📤 Sending test email to %s...\n", *to)

	_, err = svc.Users.Messages.Send("me", &gmail.Message{Raw: raw}).Context(ctx).Do()
	if err != nil {
		log.Fatalf("❌ Failed to send: %v", err)
	}

	// Persist refreshed token
	newToken, _ := tokenSource.Token()
	if newToken != nil && newToken.AccessToken != token.AccessToken {
		data, _ := json.MarshalIndent(newToken, "", "  ")
		_ = os.WriteFile(*tokenFile, data, 0600)
		fmt.Println("✅ Token refreshed and saved")
	}

	fmt.Println("✅ Email sent! Check your inbox.")
}

func buildTestMessage(from, to string) string {
	now := time.Now().Format("2006-01-02 15:04:05")

	plain := fmt.Sprintf("Health Monitor — Gmail OAuth2 Test\n\nSent at: %s\nIf you received this, OAuth2 is working!\n", now)
	html  := fmt.Sprintf(`<!DOCTYPE html><html><body style="font-family:Arial,sans-serif;max-width:500px;margin:40px auto">
<div style="background:linear-gradient(135deg,#667eea,#764ba2);color:white;padding:24px;border-radius:8px 8px 0 0;text-align:center">
<h2 style="margin:0">✅ Health Monitor</h2><p style="margin:6px 0 0;opacity:.9">Gmail OAuth2 Test</p>
</div>
<div style="border:1px solid #ddd;border-top:none;border-radius:0 0 8px 8px;padding:20px;background:#f9f9f9">
<p>Gmail OAuth2 is configured correctly!</p>
<div style="background:white;padding:10px;border-radius:4px;margin:8px 0">
  <strong>Time:</strong> %s
</div>
<div style="background:white;padding:10px;border-radius:4px;margin:8px 0">
  <strong>From:</strong> %s
</div>
</div>
<p style="text-align:center;font-size:12px;color:#888;margin-top:12px">Health Monitor — Automated Alert System</p>
</body></html>`, now, from)

	boundary := fmt.Sprintf("boundary_%d", time.Now().Unix())
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: [Health Monitor] Gmail OAuth2 Test ✅\r\n"+
			"MIME-Version: 1.0\r\nContent-Type: multipart/alternative; boundary=\"%s\"\r\n\r\n"+
			"--%s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s\r\n"+
			"--%s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s\r\n"+
			"--%s--\r\n",
		from, to, boundary, boundary, plain, boundary, html, boundary,
	)

	encoded := make([]byte, len(msg)*2)
	n := encodeBase64URL([]byte(msg), encoded)
	return string(encoded[:n])
}

func encodeBase64URL(src, dst []byte) int {
	const enc = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	n := 0
	for i := 0; i < len(src); i += 3 {
		var b [3]byte
		l := copy(b[:], src[i:])
		switch l {
		case 3:
			dst[n] = enc[b[0]>>2]; dst[n+1] = enc[(b[0]&0x3)<<4|b[1]>>4]
			dst[n+2] = enc[(b[1]&0xf)<<2|b[2]>>6]; dst[n+3] = enc[b[2]&0x3f]; n += 4
		case 2:
			dst[n] = enc[b[0]>>2]; dst[n+1] = enc[(b[0]&0x3)<<4|b[1]>>4]
			dst[n+2] = enc[(b[1]&0xf)<<2]; n += 3
		case 1:
			dst[n] = enc[b[0]>>2]; dst[n+1] = enc[(b[0]&0x3)<<4]; n += 2
		}
	}
	return n
}
