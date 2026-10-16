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
)

type IClient interface {
	TranslateNovelDetails(ctx context.Context, webpageContent string) (*models.NovelDetails, error)
	TranslateNovelChapter(ctx context.Context, novelGenres []string, webpageContent string) (*models.TranslatedChapter, error)
}

// Clients are built on first use rather than in an initialiser.
//
// Constructing them eagerly meant the server could not start without every
// provider's key even when every translation went to Bedrock, and a missing
// key panicked the process during package initialisation where no caller
// could handle it.
var (
	claudeOnce sync.Once
	claudeVal  IClient
	claudeErr  error

	bedrockOnce sync.Once
	bedrockVal  IClient
	bedrockErr  error

	litellmOnce sync.Once
	litellmVal  IClient
	litellmErr  error
)

// Get returns the client for the configured provider.
func Get(ctx context.Context) (IClient, error) {
	switch cfg := config.Get(); cfg.Provider {
	case config.ProviderClaude:
		return GetClaude(ctx)
	case config.ProviderBedrock:
		return GetBedrock(ctx)
	case config.ProviderLiteLLM:
		return GetLiteLLM(ctx)
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

// GetBedrock returns the client for any Bedrock model, through the Converse API.
func GetBedrock(ctx context.Context) (IClient, error) {
	bedrockOnce.Do(func() {
		cfg := config.Get()

		awsCfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithRegion(cfg.AWSRegion))
		if err != nil {
			bedrockErr = fmt.Errorf("load AWS config: %w", err)
			return
		}

		bedrockVal = &bedrockClientImpl{
			client:      bedrockruntime.NewFromConfig(awsCfg),
			modelID:     cfg.BedrockModelID,
			temperature: float32(cfg.TranslationTemp),
		}
	})

	return bedrockVal, bedrockErr
}

// GetLiteLLM returns the client for a LiteLLM gateway.
func GetLiteLLM(context.Context) (IClient, error) {
	litellmOnce.Do(func() {
		cfg := config.Get()

		clientConfig := openai.DefaultConfig(cfg.LiteLLMAPIKey)
		clientConfig.BaseURL = cfg.LiteLLMBaseURL

		litellmVal = newLiteLLMClient(clientConfig, cfg.LiteLLMModel, float32(cfg.TranslationTemp))
	})

	return litellmVal, litellmErr
}
