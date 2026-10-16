package llm

import (
	"context"
	"fmt"
	"log"
	"strings"

	"backend/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

// bedrockMaxTokens caps the reply length. A long chapter needs about 4,000
// output tokens, and every chat model on Bedrock that can translate a full
// chapter accepts this value.
const bedrockMaxTokens = 32000

// converser is the part of the Bedrock runtime client this provider uses.
type converser interface {
	Converse(ctx context.Context, params *bedrockruntime.ConverseInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.ConverseOutput, error)
}

// bedrockClientImpl translates through the Bedrock Converse API, which takes
// the same request for every model family on Bedrock (Anthropic, Amazon,
// Z.ai, Moonshot, DeepSeek, OpenAI and others). The claude provider instead
// sends Anthropic's native request body, which only Anthropic models accept.
type bedrockClientImpl struct {
	client converser
	// modelID is a Bedrock model id, a cross-region inference profile id, or an
	// application-inference-profile ARN. It is configuration, never a
	// compiled-in constant: an ARN embeds an AWS account id.
	modelID string
}

func (c bedrockClientImpl) TranslateNovelDetails(ctx context.Context, webpageContent string) (*models.NovelDetails, error) {
	text, err := c.converse(ctx, novelDetailsPrompt, webpageContent)
	if err != nil {
		return nil, err
	}

	var details models.NovelDetails
	if err := tryParseJSON(cleanClaudeJSON(text), &details); err != nil {
		log.Printf("bedrock %s: unparseable novel details: %v", c.modelID, err)
		return nil, fmt.Errorf("failed to unmarshal novel details: %w", err)
	}

	return &details, nil
}

func (c bedrockClientImpl) TranslateNovelChapter(ctx context.Context, novelGenres []string, webpageContent string) (*models.TranslatedChapter, error) {
	text, err := c.converse(ctx, chapterPrompt(novelGenres), webpageContent)
	if err != nil {
		return nil, err
	}

	var chapter models.TranslatedChapter
	if err := tryParseJSON(cleanClaudeJSON(text), &chapter); err != nil {
		log.Printf("bedrock %s: unparseable chapter: %v", c.modelID, err)
		return nil, fmt.Errorf("failed to unmarshal translated chapter: %w", err)
	}

	return &chapter, nil
}

// converse sends one system prompt and one user message and returns the
// reply's text. Reasoning blocks, which thinking models return alongside the
// answer, are dropped.
func (c bedrockClientImpl) converse(ctx context.Context, system, user string) (string, error) {
	out, err := c.client.Converse(ctx, &bedrockruntime.ConverseInput{
		ModelId: aws.String(c.modelID),
		System:  []types.SystemContentBlock{&types.SystemContentBlockMemberText{Value: system}},
		Messages: []types.Message{{
			Role:    types.ConversationRoleUser,
			Content: []types.ContentBlock{&types.ContentBlockMemberText{Value: user}},
		}},
		InferenceConfig: &types.InferenceConfiguration{MaxTokens: aws.Int32(bedrockMaxTokens)},
	})
	if err != nil {
		log.Printf("bedrock %s: converse: %v", c.modelID, err)
		return "", fmt.Errorf("failed to create chat completion: %w", err)
	}

	if out.StopReason == types.StopReasonMaxTokens {
		return "", fmt.Errorf("bedrock %s: reply cut off at %d tokens", c.modelID, bedrockMaxTokens)
	}

	msg, ok := out.Output.(*types.ConverseOutputMemberMessage)
	if !ok {
		return "", fmt.Errorf("bedrock %s: unexpected output type %T", c.modelID, out.Output)
	}

	var text strings.Builder
	for _, block := range msg.Value.Content {
		if t, ok := block.(*types.ContentBlockMemberText); ok {
			text.WriteString(t.Value)
		}
	}
	if text.Len() == 0 {
		return "", fmt.Errorf("bedrock %s: reply has no text", c.modelID)
	}

	return text.String(), nil
}
