//go:build ignore

// One-time OAuth2 authorization for Gmail API.
// Run this once to generate secrets/gmail-oauth-token.json.
//
// Usage:
//   go run scripts/gmail_auth.go
//   go run scripts/gmail_auth.go -credentials secrets/gmail-token.json -token secrets/gmail-oauth-token.json

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

func main() {
	credFile  := flag.String("credentials", "secrets/gmail-token.json",       "OAuth2 client credentials JSON (downloaded from Google Cloud Console)")
	tokenFile := flag.String("token",       "secrets/gmail-oauth-token.json", "Output: token file with refresh_token")
	port      := flag.String("port",        "8090",                           "Local callback server port")
	flag.Parse()

	fmt.Println("🔑 Gmail OAuth2 Authorization")
	fmt.Printf("   Credentials : %s\n", *credFile)
	fmt.Printf("   Token output: %s\n\n", *tokenFile)

	credJSON, err := os.ReadFile(*credFile)
	if err != nil {
		log.Fatalf("❌ Cannot read credentials file: %v\n\n"+
			"Download OAuth2 client JSON from:\n"+
			"  Google Cloud Console → APIs & Services → Credentials\n"+
			"  → your OAuth Client ID → Download JSON\n"+
			"Then place it at: %s\n", err, *credFile)
	}

	cfg, err := google.ConfigFromJSON(credJSON, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("❌ Invalid credentials JSON: %v", err)
	}
	cfg.RedirectURL = fmt.Sprintf("http://localhost:%s/callback", *port)

	codeCh := make(chan string, 1)
	errCh  := make(chan error, 1)

	mux := http.NewServeMux()
	srv := &http.Server{Addr: ":" + *port, Handler: mux}

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errMsg := r.URL.Query().Get("error")
			errCh <- fmt.Errorf("authorization failed: %s", errMsg)
			fmt.Fprintln(w, "<h2>❌ Authorization failed. Close this tab.</h2>")
			return
		}
		codeCh <- code
		fmt.Fprintln(w, `<html><body style="font-family:sans-serif;text-align:center;padding:60px">
<h2>✅ Authorization successful!</h2>
<p>You can close this tab and return to the terminal.</p>
</body></html>`)
	})

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("callback server: %w", err)
		}
	}()

	authURL := cfg.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	fmt.Println("🌐 Opening browser for authorization...")
	fmt.Printf("   If browser doesn't open, visit:\n   %s\n\n", authURL)
	openBrowser(authURL)

	fmt.Println("⏳ Waiting for you to authorize in the browser...")

	var code string
	select {
	case code = <-codeCh:
		fmt.Println("✅ Authorization code received")
	case err = <-errCh:
		log.Fatalf("❌ %v", err)
	}

	_ = srv.Shutdown(context.Background())

	token, err := cfg.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("❌ Failed to exchange code for token: %v", err)
	}

	if token.RefreshToken == "" {
		log.Fatal("❌ No refresh_token received.\n" +
			"Revoke app access at https://myaccount.google.com/permissions and re-run.")
	}

	if err := saveToken(*tokenFile, token); err != nil {
		log.Fatalf("❌ Failed to save token: %v", err)
	}

	rt := token.RefreshToken
	preview := rt[:8] + "..." + rt[len(rt)-4:]

	fmt.Printf("\n✅ Token saved to %s\n", *tokenFile)
	fmt.Printf("   refresh_token : %s\n", preview)
	fmt.Printf("   expiry        : %s\n", token.Expiry.Format("2006-01-02 15:04:05"))
	fmt.Println("\n🎉 Done! You can now use the gmail_oauth notifier.")
	fmt.Printf("   Next step: go run scripts/test_gmail_oauth.go\n")
}

func saveToken(path string, token *oauth2.Token) error {
	if err := os.MkdirAll("secrets", 0700); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}

func openBrowser(url string) {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "linux":
		cmd, args = "xdg-open", []string{url}
	case "darwin":
		cmd, args = "open", []string{url}
	case "windows":
		cmd, args = "rundll32", []string{"url.dll,FileProtocolHandler", url}
	default:
		return
	}
	exec.Command(cmd, args...).Start() //nolint:errcheck
}
