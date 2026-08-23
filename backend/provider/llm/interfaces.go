package llm

import (
	"context"
	"fmt"
	"sync"

	"backend/config"
	"backend/models"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	openai "github.com/sashabaranov/go-openai"
	"google.golang.org/genai"
)

type IClient interface {
	TranslateNovelDetails(ctx context.Context, webpageContent string) (*models.NovelDetails, error)
	TranslateNovelChapter(ctx context.Context, novelGenres []string, webpageContent string) (*models.TranslatedChapter, error)
}

// Clients are built on first use rather than in an initialiser.
//
// Constructing all three eagerly meant the server could not start without a
// Gemini key even when every translation went to Bedrock, and a missing key
// panicked the process during package initialisation where no caller could
// handle it.
var (
	claudeOnce sync.Once
	claudeVal  IClient
	claudeErr  error

	geminiOnce sync.Once
	geminiVal  IClient
	geminiErr  error

	openaiOnce sync.Once
	openaiVal  IClient
	openaiErr  error
)

// Get returns the client for the configured provider.
func Get(ctx context.Context) (IClient, error) {
	switch cfg := config.Get(); cfg.Provider {
	case config.ProviderClaude:
		return GetClaude(ctx)
	case config.ProviderGemini:
		return GetGemini(ctx)
	case config.ProviderOpenAI:
		return GetOpenAI(ctx)
	default:
		return nil, fmt.Errorf("unknown llm provider %q", cfg.Provider)
	}
}

// GetClaude returns the Bedrock-backed client.
func GetClaude(ctx context.Context) (IClient, error) {
	claudeOnce.Do(func() {
		cfg := config.Get()

		awsCfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithRegion(cfg.AWSRegion))
		if err != nil {
			claudeErr = fmt.Errorf("load AWS config: %w", err)
			return
		}

		claudeVal = &claudeClientImpl{
			claudeClient: bedrockruntime.NewFromConfig(awsCfg),
			modelID:      cfg.BedrockModelID,
		}
	})

	return claudeVal, claudeErr
}

// GetGemini returns the Google AI client.
func GetGemini(ctx context.Context) (IClient, error) {
	geminiOnce.Do(func() {
		cfg := config.Get()

		client, err := genai.NewClient(ctx, &genai.ClientConfig{
			APIKey:  cfg.GeminiAPIKey,
			Backend: genai.BackendGeminiAPI,
		})
		if err != nil {
			geminiErr = fmt.Errorf("create Gemini client: %w", err)
			return
		}

		geminiVal = &geminiClientImpl{geminiClient: client, model: cfg.GeminiModel}
	})

	return geminiVal, geminiErr
}

// GetOpenAI returns the OpenAI-compatible client.
func GetOpenAI(context.Context) (IClient, error) {
	openaiOnce.Do(func() {
		cfg := config.Get()

		clientConfig := openai.DefaultConfig(cfg.OpenAIAPIKey)
		if cfg.OpenAIBaseURL != "" {
			clientConfig.BaseURL = cfg.OpenAIBaseURL
		}

		openaiVal = &openaiClientImpl{
			openaiClient: openai.NewClientWithConfig(clientConfig),
			model:        cfg.OpenAIModel,
		}
	})

	return openaiVal, openaiErr
}
