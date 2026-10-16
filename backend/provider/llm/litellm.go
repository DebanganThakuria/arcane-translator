package llm

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"backend/models"

	"github.com/sashabaranov/go-openai"
)

// litellmClientImpl translates through a LiteLLM gateway, which serves many
// model families (Anthropic, OpenAI, Z.ai, Moonshot and others) behind one
// OpenAI-compatible chat completions API.
type litellmClientImpl struct {
	client *openai.Client
	model  string
	// temperature is the sampling temperature. Reasoning models such as gpt-5
	// accept only their default, so complete retries without it.
	temperature float32
}

func (c litellmClientImpl) TranslateNovelDetails(ctx context.Context, webpageContent string) (*models.NovelDetails, error) {
	text, err := c.complete(ctx, promptWithPage(novelDetailsPrompt, webpageContent))
	if err != nil {
		return nil, err
	}

	var novelDetails models.NovelDetails
	if err = tryParseJSON(cleanClaudeJSON(text), &novelDetails); err != nil {
		log.Printf("litellm %s: unparseable novel details: %v", c.model, err)
		return nil, fmt.Errorf("failed to unmarshal novel details: %w", err)
	}

	return &novelDetails, nil
}

func (c litellmClientImpl) TranslateNovelChapter(ctx context.Context, novelGenres []string, webpageContent string) (*models.TranslatedChapter, error) {
	text, err := c.complete(ctx, promptWithPage(chapterPrompt(novelGenres), webpageContent))
	if err != nil {
		return nil, err
	}

	var translatedChapter models.TranslatedChapter
	if err = tryParseJSON(cleanClaudeJSON(text), &translatedChapter); err != nil {
		log.Printf("litellm %s: unparseable chapter: %v", c.model, err)
		return nil, fmt.Errorf("failed to unmarshal translated chapter: %w", err)
	}

	return &translatedChapter, nil
}

// complete sends prompt as the only user message and returns the reply's text.
func (c litellmClientImpl) complete(ctx context.Context, prompt string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: c.temperature,
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if rejectsTemperature(err) {
		req.Temperature = 0
		resp, err = c.client.CreateChatCompletion(ctx, req)
	}
	if err != nil {
		log.Printf("litellm %s: chat completion: %v", c.model, err)
		return "", fmt.Errorf("failed to create chat completion: %w", err)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("litellm %s: reply has no choices", c.model)
	}

	log.Println(resp.Choices[0].Message.Content)

	return resp.Choices[0].Message.Content, nil
}

// rejectsTemperature reports whether err refuses the temperature parameter,
// as OpenAI's reasoning models do for any non-default value. go-openai
// refuses it before sending for model names it knows (gpt-5, o1, o3, o4);
// the API refuses it for the rest.
func rejectsTemperature(err error) bool {
	if errors.Is(err, openai.ErrReasoningModelLimitationsOther) {
		return true
	}

	var apiErr *openai.APIError
	return errors.As(err, &apiErr) &&
		apiErr.HTTPStatusCode == http.StatusBadRequest &&
		apiErr.Param != nil && *apiErr.Param == "temperature"
}
