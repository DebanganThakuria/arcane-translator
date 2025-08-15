package llm

import (
	"context"
	"os"

	"backend/models"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/genai"
)

var (
	geminiClient IClient
	openaiClient IClient
	claudeClient IClient
)

type IClient interface {
	TranslateNovelDetails(ctx context.Context, webpageContent string) (*models.NovelDetails, error)
	TranslateNovelChapter(ctx context.Context, novelGenres []string, webpageContent string) (*models.TranslatedChapter, error)
}

func init() {
	genaiClient, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		panic("Failed to create Gemini client: " + err.Error())
	}

	geminiClient = &geminiClientImpl{
		geminiClient: genaiClient,
	}

	config := openai.DefaultConfig(os.Getenv("OPENAI_API_KEY"))
	config.BaseURL = os.Getenv("OPENAI_API_BASE_URL")

	openaiClient = &openaiClientImpl{
		openaiClient: openai.NewClientWithConfig(config),
	}

	cfg, _ := awsConfig.LoadDefaultConfig(context.Background(), awsConfig.WithRegion("us-east-1"))
	client := bedrockruntime.NewFromConfig(cfg)

	claudeClient = &claudeClientImpl{
		claudeClient: client,
	}
}

func GetGemini() IClient {
	return geminiClient
}

func GetOpenAI() IClient {
	return openaiClient
}

func GetClaude() IClient {
	return claudeClient
}
