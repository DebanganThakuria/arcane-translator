package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"backend/models"

	"github.com/sashabaranov/go-openai"
)

// litellmMaxTokens caps the reply length. The gateway's own default is 4,096
// output tokens, which ends a long chapter mid-JSON. A normal chapter needs
// about 5,000; a "ten thousand words" chapter of 8,700 Chinese characters
// needed about 17,000, so this leaves room for one twice that long. Bedrock
// models take the same cap.
const litellmMaxTokens = 32000

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

// newLiteLLMClient builds the client for model. For models that think by
// default but can be told not to, every request turns thinking off: a
// translation does not need it, and on Claude Sonnet 5.5 it added about 500
// output tokens (16%) to a chapter with no change in quality.
func newLiteLLMClient(cfg openai.ClientConfig, model string, temperature float32) *litellmClientImpl {
	if fields := cheapModeFields(model); fields != nil {
		cfg.HTTPClient = extraFieldsDoer{next: cfg.HTTPClient, fields: fields}
	}

	return &litellmClientImpl{
		client:      openai.NewClientWithConfig(cfg),
		model:       model,
		temperature: temperature,
	}
}

// cheapModeFields returns the request fields that turn thinking off and lower
// effort for model, or nil when there are none.
//
// Claude Sonnet 5.5 rejects thinking "disabled" and takes "between_tools"
// instead; every other model rejects "between_tools", so it must never be sent
// to them. Claude Opus 5.5 cannot turn thinking off. With thinking off, medium
// effort translated a chapter exactly as fully as the default high, with
// slightly fewer tokens; low wrote noticeably tighter prose for little saving.
func cheapModeFields(model string) map[string]any {
	if strings.Contains(strings.ReplaceAll(model, ".", "-"), "sonnet-5-5") {
		return map[string]any{
			"thinking":      map[string]any{"type": "between_tools"},
			"output_config": map[string]any{"effort": "medium"},
		}
	}
	return nil
}

// extraFieldsDoer adds fields to every chat completion request body.
// go-openai has no field for provider-specific parameters such as Anthropic's
// thinking setting, which the LiteLLM gateway passes through to the model.
type extraFieldsDoer struct {
	next   openai.HTTPDoer
	fields map[string]any
}

func (d extraFieldsDoer) Do(req *http.Request) (*http.Response, error) {
	if req.Body == nil || !strings.HasSuffix(req.URL.Path, "/chat/completions") {
		return d.next.Do(req)
	}

	var body map[string]any
	dec := json.NewDecoder(req.Body)
	dec.UseNumber()
	err := dec.Decode(&body)
	req.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("read chat completion request: %w", err)
	}

	for k, v := range d.fields {
		body[k] = v
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("write chat completion request: %w", err)
	}

	req.Body = io.NopCloser(bytes.NewReader(raw))
	req.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(raw)), nil }
	req.ContentLength = int64(len(raw))
	return d.next.Do(req)
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
		Temperature:         c.temperature,
		MaxCompletionTokens: litellmMaxTokens,
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
	if resp.Choices[0].FinishReason == openai.FinishReasonLength {
		log.Printf("litellm %s: reply cut off at %d tokens", c.model, litellmMaxTokens)
		return "", fmt.Errorf("litellm %s: reply cut off at %d tokens", c.model, litellmMaxTokens)
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
