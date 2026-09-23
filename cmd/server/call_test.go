package main

import "testing"

func TestBuildCallConfig(t *testing.T) {
	stored := map[string]interface{}{
		"phones":           []interface{}{"79000000000"},
		"voice":            "GOLOSM",
		"message_template": "Сервис {{.TargetName}} не отвечает",
		"api_key_env":      "CUSTOM_KEY",
		"min_failures":     3,
	}

	t.Run("saved settings are kept", func(t *testing.T) {
		cfg, err := buildCallConfig(stored, "", "", "")
		if err != nil {
			t.Fatalf("buildCallConfig failed: %v", err)
		}
		if cfg["message_template"] != stored["message_template"] {
			t.Errorf("Expected the saved template, got %v", cfg["message_template"])
		}
		if cfg["api_key_env"] != "CUSTOM_KEY" {
			t.Errorf("Expected the saved key variable, got %v", cfg["api_key_env"])
		}
		if cfg["voice"] != "GOLOSM" {
			t.Errorf("Expected the saved voice, got %v", cfg["voice"])
		}
	})

	t.Run("flags override", func(t *testing.T) {
		cfg, err := buildCallConfig(stored, "79001234567, 79007654321", "Проверка связи", "GOLOS")
		if err != nil {
			t.Fatalf("buildCallConfig failed: %v", err)
		}

		phones, ok := cfg["phones"].([]interface{})
		if !ok || len(phones) != 2 || phones[0] != "79001234567" || phones[1] != "79007654321" {
			t.Fatalf("Expected both numbers from -to, got %v", cfg["phones"])
		}
		if cfg["message_template"] != "Проверка связи" {
			t.Errorf("Expected -text to win, got %v", cfg["message_template"])
		}
		if cfg["voice"] != "GOLOS" {
			t.Errorf("Expected -voice to win, got %v", cfg["voice"])
		}

		// The saved config must not be mutated by a call.
		if len(stored["phones"].([]interface{})) != 1 {
			t.Error("Expected the stored config to stay untouched")
		}
	})

	t.Run("no phone anywhere", func(t *testing.T) {
		if _, err := buildCallConfig(map[string]interface{}{}, "", "hi", ""); err == nil {
			t.Fatal("Expected an error without a phone number")
		}
	})

	t.Run("default voice", func(t *testing.T) {
		cfg, err := buildCallConfig(map[string]interface{}{}, "79000000000", "", "")
		if err != nil {
			t.Fatalf("buildCallConfig failed: %v", err)
		}
		if cfg["voice"] != "GOLOS" {
			t.Errorf("Expected GOLOS by default, got %v", cfg["voice"])
		}
	})
}
