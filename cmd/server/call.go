package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"health-monitor/internal/domain"
	"health-monitor/internal/notifier"
	"health-monitor/internal/storage"
	"health-monitor/pkg/config"
)

// Exit codes, so a wrapper script can tell the outcomes apart.
const (
	callExitOK        = 0
	callExitUsage     = 1
	callExitRejected  = 2
	callExitUnknown   = 3
	callPlaceTimeout  = 60 * time.Second
	callDefaultTarget = "Test target"
)

// runCallCommand places a single voice call by hand. It exists because the
// SMSPILOT key lives only in the server's environment: the dashboard has no
// test button on purpose (a call costs money and rings a real person), so the
// check belongs on the machine that holds the key.
//
// It touches neither incidents nor the notification journal.
func runCallCommand(args []string) int {
	fs := flag.NewFlagSet("call", flag.ContinueOnError)

	to := fs.String("to", "", "Phone number(s) to call, comma separated (79XXXXXXXXX)")
	text := fs.String("text", "", "What to say; may use {{.TargetName}} (default: the notifier's template)")
	voice := fs.String("voice", "", "GOLOS (female) or GOLOSM (male)")
	notifierID := fs.String("notifier", "", "Place the call with a saved notifier's settings (template, voice, phones)")
	configPath := fs.String("config", defaultCallConfigPath(), "Path to configuration file, used to find the database")
	targetName := fs.String("target", callDefaultTarget, "Target name substituted into the message")
	dryRun := fs.Bool("dry-run", false, "Print what would be said and dialled, place no call")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Place one voice call, for checking the SMSPILOT setup.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n  health-monitor call -to 79001234567 [-text \"...\"] [-notifier oncall-phone] [-dry-run]\n\n")
		fmt.Fprintf(os.Stderr, "The API key is read from SMSPILOT_API_KEY (or the notifier's api_key_env).\n\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return callExitUsage
	}

	stored := map[string]interface{}{}
	if *notifierID != "" {
		loaded, err := loadNotifierConfig(*configPath, *notifierID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load notifier %q: %v\n", *notifierID, err)
			return callExitUsage
		}
		stored = loaded
	}

	callConfig, err := buildCallConfig(stored, *to, *text, *voice)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		fs.Usage()
		return callExitUsage
	}

	voiceNotifier, err := notifier.NewSMSPilotNotifier(callConfig, zerolog.Nop())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid voice configuration: %v\n", err)
		return callExitUsage
	}

	alert := &domain.Alert{
		ID:         "manual-call",
		TargetID:   "manual-call",
		TargetName: *targetName,
		Type:       domain.AlertTypeDown,
		Severity:   domain.AlertSeverityCritical,
		Message:    fmt.Sprintf("Target %s is DOWN", *targetName),
		CreatedAt:  time.Now(),
		// Sample values, so a preview of a template that uses them shows what
		// would actually be said instead of "<no value>".
		Description: "Unexpected status code: got 502, expected 200",
		Metadata: map[string]interface{}{
			"incident_id":   int64(42),
			"failure_count": 3,
			"target_type":   "http",
		},
	}

	message, err := voiceNotifier.RenderMessage(alert)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to render the message: %v\n", err)
		return callExitUsage
	}

	phones := callConfig["phones"].([]interface{})

	fmt.Printf("Calling:  %v\n", phones)
	fmt.Printf("Voice:    %v\n", callConfig["voice"])
	fmt.Printf("Message:  %s\n", message)

	if *dryRun {
		fmt.Println("\nDry run: no call was placed.")
		return callExitOK
	}

	ctx, cancel := context.WithTimeout(context.Background(), callPlaceTimeout)
	defer cancel()

	receipt, err := voiceNotifier.NotifyWithReceipt(ctx, alert)

	switch {
	case err == nil:
		fmt.Printf("\nAccepted. server_id=%s status=%s cost=%s\n", receipt.ExternalID, receipt.ExternalStatus, receipt.Cost)
		fmt.Println("Status 2 means the call was delivered, not that anyone picked up or acted on it.")
		return callExitOK

	case errors.Is(err, domain.ErrUnknownResult):
		fmt.Fprintf(os.Stderr, "\nResult unknown: %v\n", err)
		fmt.Fprintln(os.Stderr, "Do not simply run this again: the call may already have been placed.")
		return callExitUnknown

	default:
		fmt.Fprintf(os.Stderr, "\nRejected: %v\n", err)
		return callExitRejected
	}
}

// defaultCallConfigPath picks the config the server itself is running with.
// In the container it lives at /configs/example.yaml (see the Dockerfile CMD),
// in a checkout it is relative — so `-notifier` works in both without a flag.
func defaultCallConfigPath() string {
	for _, path := range []string{"configs/example.yaml", "/configs/example.yaml"} {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return "configs/example.yaml"
}

// buildCallConfig merges the saved notifier settings with the flags, so a call
// can exercise the exact configuration an incident would use.
func buildCallConfig(stored map[string]interface{}, to, text, voice string) (map[string]interface{}, error) {
	callConfig := make(map[string]interface{}, len(stored)+3)
	for k, v := range stored {
		callConfig[k] = v
	}

	if to != "" {
		phones := make([]interface{}, 0, 2)
		for _, phone := range strings.Split(to, ",") {
			if phone = strings.TrimSpace(phone); phone != "" {
				phones = append(phones, phone)
			}
		}
		callConfig["phones"] = phones
	}

	if phones, ok := callConfig["phones"].([]interface{}); !ok || len(phones) == 0 {
		return nil, fmt.Errorf("-to is required (or -notifier of a channel that has phones)")
	}

	if text != "" {
		callConfig["message_template"] = text
	}

	if voice != "" {
		callConfig["voice"] = voice
	}

	if _, ok := callConfig["voice"]; !ok {
		callConfig["voice"] = "GOLOS"
	}

	return callConfig, nil
}

// loadNotifierConfig reads a saved voice notifier from the database.
func loadNotifierConfig(configPath, notifierID string) (map[string]interface{}, error) {
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	db, err := storage.New(cfg.Database, zerolog.Nop())
	if err != nil {
		return nil, fmt.Errorf("failed to open the database: %w", err)
	}
	defer db.Close()

	stored, err := storage.NewNotifierRepository(db.DB()).Get(context.Background(), notifierID)
	if err != nil {
		return nil, err
	}

	if stored.Type != domain.NotifierTypeSMSPilot {
		return nil, fmt.Errorf("notifier %q is of type %q, not %q", notifierID, stored.Type, domain.NotifierTypeSMSPilot)
	}

	return stored.Config, nil
}
