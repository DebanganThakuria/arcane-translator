package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateRequiresProviderSettings(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{"claude with a model", Config{Provider: ProviderClaude, BedrockModelID: "model"}, false},
		{"claude without a model", Config{Provider: ProviderClaude}, true},
		{"gemini with a key", Config{Provider: ProviderGemini, GeminiAPIKey: "key"}, false},
		{"gemini without a key", Config{Provider: ProviderGemini}, true},
		{"openai with a key", Config{Provider: ProviderOpenAI, OpenAIAPIKey: "key"}, false},
		{"openai without a key", Config{Provider: ProviderOpenAI}, true},
		{"unknown provider", Config{Provider: "llama"}, true},
		{"empty provider", Config{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr && err == nil {
				t.Error("expected an error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestFromEnvReadsProviderSettings(t *testing.T) {
	t.Setenv("ARCANE_LLM_PROVIDER", "GEMINI")
	t.Setenv("GEMINI_API_KEY", "abc123")
	t.Setenv("ARCANE_PORT", "9999")
	t.Setenv("ARCANE_DB_PATH", "/tmp/test.db")

	cfg := fromEnv()

	if cfg.Provider != ProviderGemini {
		t.Errorf("provider = %q, want gemini (case-insensitive)", cfg.Provider)
	}
	if cfg.GeminiAPIKey != "abc123" {
		t.Errorf("GeminiAPIKey = %q", cfg.GeminiAPIKey)
	}
	if cfg.Port != 9999 {
		t.Errorf("Port = %d, want 9999", cfg.Port)
	}
	if cfg.DBPath != "/tmp/test.db" {
		t.Errorf("DBPath = %q", cfg.DBPath)
	}
}

func TestFromEnvDefaults(t *testing.T) {
	for _, key := range []string{
		"ARCANE_LLM_PROVIDER", "ARCANE_PORT", "AWS_REGION", "ARCANE_GEMINI_MODEL",
	} {
		t.Setenv(key, "")
	}
	t.Setenv("ARCANE_DB_PATH", "/tmp/test.db")

	cfg := fromEnv()

	if cfg.Provider != ProviderClaude {
		t.Errorf("default provider = %q, want claude", cfg.Provider)
	}
	if cfg.Port != 8088 {
		t.Errorf("default port = %d, want 8088", cfg.Port)
	}
	if cfg.AWSRegion != "us-east-1" {
		t.Errorf("default region = %q", cfg.AWSRegion)
	}
	if cfg.GeminiModel == "" {
		t.Error("gemini model default is empty")
	}
}

// An invalid port must not take the server to port 0.
func TestFromEnvIgnoresBadPort(t *testing.T) {
	t.Setenv("ARCANE_DB_PATH", "/tmp/test.db")

	for _, bad := range []string{"not-a-number", "-1", "0"} {
		t.Setenv("ARCANE_PORT", bad)
		if got := fromEnv().Port; got != 8088 {
			t.Errorf("ARCANE_PORT=%q gave port %d, want the 8088 default", bad, got)
		}
	}
}

// A real environment variable must win over the .env file, so a deployment is
// never silently overridden by a file left in the working directory.
func TestLoadDotEnvDoesNotOverrideEnvironment(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	contents := "ARCANE_TEST_PRESET=from-file\nARCANE_TEST_UNSET=from-file\n# comment\n\n"
	if err := os.WriteFile(envFile, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Chdir(dir)
	t.Setenv("ARCANE_TEST_PRESET", "from-environment")
	os.Unsetenv("ARCANE_TEST_UNSET")
	t.Cleanup(func() { os.Unsetenv("ARCANE_TEST_UNSET") })

	loadDotEnv()

	if got := os.Getenv("ARCANE_TEST_PRESET"); got != "from-environment" {
		t.Errorf("preset variable = %q, want the environment to win", got)
	}
	if got := os.Getenv("ARCANE_TEST_UNSET"); got != "from-file" {
		t.Errorf("unset variable = %q, want the file value", got)
	}
}
