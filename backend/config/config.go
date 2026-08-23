// Package config resolves runtime settings from the environment.
//
// Nothing secret is compiled into the binary: model identifiers, API keys and
// the AWS account behind a Bedrock inference profile all arrive as environment
// variables, so the source tree can be published as-is.
package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// Provider names the translation backend in use.
type Provider string

const (
	ProviderClaude Provider = "claude"
	ProviderGemini Provider = "gemini"
	ProviderOpenAI Provider = "openai"
)

// Config is the resolved runtime configuration.
type Config struct {
	Port   int
	DBPath string
	// WebDir holds the built frontend. Empty in development, where Vite serves
	// it; set for a packaged install so one process serves everything.
	WebDir string

	Provider Provider

	AWSRegion       string
	BedrockModelID  string
	GeminiAPIKey    string
	GeminiModel     string
	OpenAIAPIKey    string
	OpenAIBaseURL   string
	OpenAIModel     string
	TranslationTemp float64
}

var (
	once     sync.Once
	resolved *Config
)

// Get returns the process configuration, loading it once.
func Get() *Config {
	once.Do(func() {
		loadDotEnv()
		resolved = fromEnv()
	})
	return resolved
}

// loadDotEnv reads KEY=VALUE lines from a .env file next to the binary or in
// the working directory. Real environment variables always win, so this is a
// convenience for local development and never overrides a deployment.
func loadDotEnv() {
	for _, path := range dotEnvCandidates() {
		file, err := os.Open(path)
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			key, value, found := strings.Cut(line, "=")
			if !found {
				continue
			}

			key = strings.TrimSpace(key)
			value = strings.Trim(strings.TrimSpace(value), `"'`)
			if _, exists := os.LookupEnv(key); !exists {
				_ = os.Setenv(key, value)
			}
		}
		file.Close()
		return
	}
}

func dotEnvCandidates() []string {
	candidates := []string{".env"}

	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates,
			filepath.Join(cwd, ".env"),
			filepath.Join(filepath.Dir(cwd), ".env"),
		)
	}

	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(home, ".arcane-translator", ".env"))
	}

	return candidates
}

func fromEnv() *Config {
	return &Config{
		Port:   envInt("ARCANE_PORT", 8088),
		DBPath: resolveDBPath(),
		WebDir: resolveWebDir(),

		Provider: Provider(strings.ToLower(envString("ARCANE_LLM_PROVIDER", string(ProviderClaude)))),

		AWSRegion:       envString("AWS_REGION", "us-east-1"),
		BedrockModelID:  os.Getenv("ARCANE_BEDROCK_MODEL_ID"),
		GeminiAPIKey:    os.Getenv("GEMINI_API_KEY"),
		GeminiModel:     envString("ARCANE_GEMINI_MODEL", "gemini-2.5-flash"),
		OpenAIAPIKey:    os.Getenv("OPENAI_API_KEY"),
		OpenAIBaseURL:   os.Getenv("OPENAI_API_BASE_URL"),
		OpenAIModel:     envString("ARCANE_OPENAI_MODEL", "gpt-5"),
		TranslationTemp: envFloat("ARCANE_TRANSLATION_TEMPERATURE", 0.2),
	}
}

// resolveDBPath decides where the library lives. It always returns a writable
// location: an explicit setting, a repository checkout, or the user's data
// directory. It never falls back to a path relative to the working directory,
// which outside a checkout resolved to /data and failed to open.
func resolveDBPath() string {
	if explicit := os.Getenv("ARCANE_DB_PATH"); explicit != "" {
		return explicit
	}

	if dir := os.Getenv("ARCANE_DATA_DIR"); dir != "" {
		return filepath.Join(dir, "data.db")
	}

	// Running from a checkout keeps its own data/ directory, so a developer's
	// library is not silently moved to the packaged location.
	if root := repoRoot(); root != "" {
		dataDir := filepath.Join(root, "data")
		if err := os.MkdirAll(dataDir, 0o755); err == nil {
			return filepath.Join(dataDir, "data.db")
		}
	}

	if home, err := os.UserHomeDir(); err == nil {
		dir := filepath.Join(home, ".arcane-translator")
		if err := os.MkdirAll(dir, 0o755); err == nil {
			return filepath.Join(dir, "data.db")
		}
	}

	return filepath.Join(os.TempDir(), "arcane-translator.db")
}

// repoRoot walks up from the working directory looking for backend/go.mod,
// returning the repository root, or "" when not running inside a checkout.
func repoRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for range 6 {
		if _, err := os.Stat(filepath.Join(dir, "backend", "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}

	return ""
}

// resolveWebDir finds the built frontend: an explicit path, a sibling of the
// binary as laid out by the Homebrew formula, or the repository's web/build.
func resolveWebDir() string {
	if explicit := os.Getenv("ARCANE_WEB_DIR"); explicit != "" {
		return explicit
	}

	candidates := []string{}
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Join(dir, "web"),
			filepath.Join(dir, "..", "share", "arcane-translator", "web"),
		)
	}
	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates,
			filepath.Join(cwd, "web", "build"),
			filepath.Join(cwd, "..", "web", "build"),
		)
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(filepath.Join(candidate, "index.html")); err == nil {
			return candidate
		}
	}

	return ""
}

// Validate reports whether the selected provider has what it needs. It is
// checked at startup so a misconfigured server fails with a clear message
// instead of panicking deep inside a package initialiser.
func (c *Config) Validate() error {
	switch c.Provider {
	case ProviderClaude:
		if c.BedrockModelID == "" {
			return fmt.Errorf(
				"ARCANE_BEDROCK_MODEL_ID is required for the claude provider " +
					"(a Bedrock model id or an application-inference-profile ARN)")
		}
	case ProviderGemini:
		if c.GeminiAPIKey == "" {
			return fmt.Errorf("GEMINI_API_KEY is required for the gemini provider")
		}
	case ProviderOpenAI:
		if c.OpenAIAPIKey == "" {
			return fmt.Errorf("OPENAI_API_KEY is required for the openai provider")
		}
	default:
		return fmt.Errorf(
			"unknown ARCANE_LLM_PROVIDER %q, expected claude, gemini or openai", c.Provider)
	}

	return nil
}

func envString(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if value, err := strconv.Atoi(os.Getenv(key)); err == nil && value > 0 {
		return value
	}
	return fallback
}

func envFloat(key string, fallback float64) float64 {
	if value, err := strconv.ParseFloat(os.Getenv(key), 64); err == nil {
		return value
	}
	return fallback
}
