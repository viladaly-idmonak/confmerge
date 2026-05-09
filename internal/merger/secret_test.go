package merger

import (
	"testing"
)

func TestSecretMasker_NoPatterns(t *testing.T) {
	masker, err := NewSecretMasker(nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := map[string]any{"host": "localhost", "port": 5432}
	result := masker.Mask(input)
	if result["host"] != "localhost" || result["port"] != 5432 {
		t.Errorf("expected unchanged map, got %v", result)
	}
}

func TestSecretMasker_SimplePattern(t *testing.T) {
	masker, err := NewSecretMasker([]string{"*password*"}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := map[string]any{
		"db_password": "s3cr3t",
		"username":    "admin",
	}
	result := masker.Mask(input)
	if result["db_password"] != "***REDACTED***" {
		t.Errorf("expected redacted password, got %v", result["db_password"])
	}
	if result["username"] != "admin" {
		t.Errorf("expected unchanged username, got %v", result["username"])
	}
}

func TestSecretMasker_CaseInsensitive(t *testing.T) {
	masker, err := NewSecretMasker([]string{"*SECRET*"}, "[hidden]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := map[string]any{"api_secret": "abc123", "name": "app"}
	result := masker.Mask(input)
	if result["api_secret"] != "[hidden]" {
		t.Errorf("expected [hidden], got %v", result["api_secret"])
	}
}

func TestSecretMasker_NestedMap(t *testing.T) {
	masker, err := NewSecretMasker([]string{"*token*", "*key*"}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := map[string]any{
		"service": map[string]any{
			"api_key": "xyz",
			"url":     "https://example.com",
		},
		"auth_token": "tok123",
	}
	result := masker.Mask(input)
	if result["auth_token"] != "***REDACTED***" {
		t.Errorf("expected redacted token")
	}
	svc, ok := result["service"].(map[string]any)
	if !ok {
		t.Fatal("expected nested map")
	}
	if svc["api_key"] != "***REDACTED***" {
		t.Errorf("expected redacted api_key, got %v", svc["api_key"])
	}
	if svc["url"] != "https://example.com" {
		t.Errorf("expected unchanged url")
	}
}

func TestSecretMasker_InvalidPattern(t *testing.T) {
	_, err := NewSecretMasker([]string{"[invalid"}, "")
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}
