package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"backend/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type claudeClientImpl struct {
	claudeClient *bedrockruntime.Client
}

type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnthropicMessagesRequest struct {
	AnthropicVersion string             `json:"anthropic_version,omitempty"`
	Messages         []AnthropicMessage `json:"messages"`
	MaxTokens        int                `json:"max_tokens"`
	Temperature      float64            `json:"temperature,omitempty"`
	TopP             float64            `json:"top_p,omitempty"`
	System           string             `json:"system,omitempty"`
}

type AnthropicResponse struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Role    string `json:"role"`
	Model   string `json:"model"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	StopReason   string      `json:"stop_reason"`
	StopSequence interface{} `json:"stop_sequence"`
	Usage        struct {
		InputTokens              int `json:"input_tokens"`
		CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
		CacheReadInputTokens     int `json:"cache_read_input_tokens"`
		OutputTokens             int `json:"output_tokens"`
	} `json:"usage"`
}

func (c claudeClientImpl) TranslateNovelDetails(ctx context.Context, webpageContent string) (*models.NovelDetails, error) {
	prompt := `
	You are a professional translator for webnovels.
	Please extract and translate information from this novel page content.

	Return ONLY a valid JSON object with this exact structure:
	{
		"novel_title_original": "Original title in the source language",
		"novel_title_translated": "Translated title in English",
		"novel_summary_translated": "Translated summary in English in HTML format with paragraph tags. Please ensure that the summary has valid HTML tags for rendering on the frontend.",
		"possible_novel_genres": ["Genre1", "Genre2", ...],
		"number_of_chapters": "total number of chapters in original as integer",
		"status": "Ongoing" or "Completed" or "Unknown"
	}

	Do not include any commentary, explanation, or preamble. Only return the valid JSON object. It should be correctly and directly marshallable into a go struct.
	Do NOT wrap the JSON object in any Markdown code block. Return only the raw JSON object, with no extra formatting.`

	modelOrProfile := ""

	req := AnthropicMessagesRequest{
		AnthropicVersion: "bedrock-2023-05-31",
		Messages: []AnthropicMessage{
			{
				Role:    "user",
				Content: webpageContent,
			},
		},
		System:    prompt,
		MaxTokens: 50000,
	}

	body, err := json.Marshal(req)
	if err != nil {
		log.Println(ctx, "Marshal err: %v", err)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	out, err := c.claudeClient.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(modelOrProfile), // application inference profile ARN
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		log.Println(ctx, "CreateChatCompletion err: %v", err)
		return nil, fmt.Errorf("failed to create chat completion: %w", err)
	}

	log.Println(string(out.Body))

	var response AnthropicResponse
	if err = json.Unmarshal(out.Body, &response); err != nil {
		log.Println(ctx, "Unmarshal err: %v", err)
		return nil, err
	}

	var novelDetails models.NovelDetails
	if err = json.Unmarshal([]byte(response.Content[0].Text), &novelDetails); err != nil {
		log.Println(ctx, "Unmarshal err: %v", err)
		return nil, fmt.Errorf("failed to unmarshal novel details: %w", err)
	}

	return &novelDetails, nil
}

func (c claudeClientImpl) TranslateNovelChapter(ctx context.Context, novelGenres []string, webpageContent string) (*models.TranslatedChapter, error) {
	prompt := `
        You are the best webnovel translator and editor, capable of producing the highest quality work.
		Your task is translating and polishing the following Webnovel chapter into flawless English, ensuring perfect grammar and language. Translate all original language object names, including places, abilities, techniques, and other cultural references, into English.
		Your goal is to craft an engaging English version while preserving the original content. The author of the original text may have limited writing skills, so take the liberty to polish the content in your translation wherever you judge necessary. Pay special attention to the dialogues, ensuring they flow smoothly and sound lifelike. Do not add anything extra anywhere from your end.
		Finally, you mustn't lose any content from the original during the translation process. 
		I trust you to provide the best possible results. Please translate the full chapter as per these guidelines.

        Currently known novel genres: ` + models.GenresToString(novelGenres) + `

        Return ONLY a valid JSON object with this exact structure:
        {{
            "translated_chapter_title": "The translated title of the chapter.",
			"original_chapter_title": "The original title in the source language.",
            "translated_chapter_contents": "The translated content of the full chapter in HTML format with paragraph tags. Please ensure that the chapter content has valid HTML tags for rendering on the frontend. And, most importantly, ensure that the full chapter content is included in the response.",
            "possible_new_genres": "Array of any genres you detect that aren't already known (leave empty if none). Only check for standard genres. Otherwise, this list grows exponentially. Response should be like this ["Genre1", "Genre2", ...]"
        }}
	
		Additional instructions. Please follow these instructions carefully:
		- Please terminate the translation and return when the chapter ends.
		- Do not miss any of the fields. This is very important. 
		- Do not include any commentary, explanation, or preamble. Only return the JSON object. It should be correctly and directly marshallable into a go struct.
	    - Do NOT wrap the JSON object in any Markdown code block. Return only the raw JSON object, with no extra formatting.
	`
	modelOrProfile := ""

	req := AnthropicMessagesRequest{
		AnthropicVersion: "bedrock-2023-05-31",
		Messages: []AnthropicMessage{
			{
				Role:    "user",
				Content: webpageContent,
			},
		},
		System:    prompt,
		MaxTokens: 50000,
	}

	body, _ := json.Marshal(req)

	out, err := c.claudeClient.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(modelOrProfile), // application inference profile ARN
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		log.Println(ctx, "CreateChatCompletion err: %v", err)
		return nil, fmt.Errorf("failed to create chat completion: %w", err)
	}

	log.Println(string(out.Body))

	var response AnthropicResponse
	if err = json.Unmarshal(out.Body, &response); err != nil {
		log.Println(ctx, "Unmarshal err: %v", err)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	var translatedChapter models.TranslatedChapter
	if err = json.Unmarshal([]byte(response.Content[0].Text), &translatedChapter); err != nil {
		log.Println(ctx, "Unmarshal err: %v", err)
		return nil, err
	}

	return &translatedChapter, nil
}
