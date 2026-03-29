package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

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

// cleanClaudeJSON trims whitespace, removes markdown/code block markers, and attempts to extract the first JSON object from the string.
func cleanClaudeJSON(s string) string {
	s = strings.TrimSpace(s)
	// Remove Markdown code block markers if present
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	s = strings.TrimSpace(s)
	// Try to extract the first JSON object using regex
	re := regexp.MustCompile(`(?s){.*}`)
	match := re.FindString(s)
	if match != "" {
		return match
	}
	return s
}

// fixUnescapedQuotesInJSON attempts to fix unescaped quotes within JSON string values.
// This is needed when LLMs produce JSON with dialogue containing unescaped quotes.
func fixUnescapedQuotesInJSON(s string) string {
	var result strings.Builder
	inString := false
	escaped := false
	stringStart := -1

	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		c := runes[i]

		if escaped {
			result.WriteRune(c)
			escaped = false
			continue
		}

		if c == '\\' {
			result.WriteRune(c)
			escaped = true
			continue
		}

		if c == '"' {
			if !inString {
				// Starting a string
				inString = true
				stringStart = i
				result.WriteRune(c)
			} else {
				// Potentially ending a string - but we need to check if this is a legitimate end
				// Look ahead to see if next non-whitespace char is :, ,, }, or ]
				isLegitEnd := false
				for j := i + 1; j < len(runes); j++ {
					nextChar := runes[j]
					if nextChar == ' ' || nextChar == '\t' || nextChar == '\n' || nextChar == '\r' {
						continue
					}
					if nextChar == ':' || nextChar == ',' || nextChar == '}' || nextChar == ']' {
						isLegitEnd = true
					}
					break
				}

				// Also check if this might be starting a new key (look for pattern: "key":)
				// by checking if there's text before the quote that looks like a value ending
				if !isLegitEnd {
					// Look for pattern where quote is followed by text then colon (new key)
					// This handles cases like: value" , "newKey":
					for j := i + 1; j < len(runes); j++ {
						nextChar := runes[j]
						if nextChar == ' ' || nextChar == '\t' || nextChar == '\n' || nextChar == '\r' {
							continue
						}
						// If we see a letter/underscore, this might be continuing content or a new key
						// Check further ahead for colon to determine
						if (nextChar >= 'a' && nextChar <= 'z') || (nextChar >= 'A' && nextChar <= 'Z') || nextChar == '_' {
							// Scan for quote-colon pattern to determine if this is a new key
							foundQuoteColon := false
							for k := j; k < len(runes); k++ {
								if runes[k] == '"' {
									// Check if pattern is "key":
									for l := k + 1; l < len(runes); l++ {
										if runes[l] == ' ' || runes[l] == '\t' {
											continue
										}
										if runes[l] == ':' {
											foundQuoteColon = true
										}
										break
									}
									break
								}
								if runes[k] == ',' || runes[k] == '}' || runes[k] == ']' {
									break
								}
							}
							if foundQuoteColon && stringStart > 0 {
								// This is likely dialogue - the quote should be escaped
								result.WriteString("\\\"")
								inString = true // stay in string
								continue
							}
						}
						break
					}
				}

				if isLegitEnd {
					inString = false
					stringStart = -1
				} else {
					// This quote appears to be inside dialogue - escape it
					result.WriteRune('\\')
				}
				result.WriteRune(c)
			}
		} else {
			result.WriteRune(c)
		}
	}

	return result.String()
}

// tryParseJSON attempts to parse JSON, and if it fails, tries to fix common issues
func tryParseJSON[T any](jsonStr string, target *T) error {
	// First, try parsing as-is
	err := json.Unmarshal([]byte(jsonStr), target)
	if err == nil {
		return nil
	}

	// If that fails, try fixing unescaped quotes
	fixedJSON := fixUnescapedQuotesInJSON(jsonStr)
	err2 := json.Unmarshal([]byte(fixedJSON), target)
	if err2 == nil {
		return nil
	}

	// Return the original error if both attempts fail
	return fmt.Errorf("failed to parse JSON (original error: %v, fixed attempt error: %v)", err, err2)
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

	modelOrProfile := "arn:aws:bedrock:us-east-1:101860328116:application-inference-profile/5o3okv42vqwm"

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
	jsonStr := cleanClaudeJSON(response.Content[0].Text)
	if err = tryParseJSON(jsonStr, &novelDetails); err != nil {
		log.Println(ctx, "Unmarshal err: %v. Raw response: %s", err, response.Content[0].Text)
		return nil, fmt.Errorf("failed to unmarshal novel details: %w", err)
	}

	return &novelDetails, nil
}

func (c claudeClientImpl) TranslateNovelChapter(ctx context.Context, novelGenres []string, webpageContent string) (*models.TranslatedChapter, error) {
	prompt := `
        You are the best webnovel translator and editor, capable of producing the highest quality work.
		Your task is translating and polishing the following Webnovel chapter into flawless English, ensuring perfect grammar and language. Translate all original language object names, including places, abilities, techniques, and other cultural references, into English.
		Your goal is to craft an engaging English version while preserving the original content. The author of the original text may have limited writing skills, so take the liberty to polish the content in your translation to ensure that the sentences and paragraphs flow smoothly. Pay special attention to the dialogues, ensuring they flow smoothly and sound lifelike.
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
	modelOrProfile := "arn:aws:bedrock:us-east-1:101860328116:application-inference-profile/5o3okv42vqwm"

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
	jsonStr := cleanClaudeJSON(response.Content[0].Text)
	if err = tryParseJSON(jsonStr, &translatedChapter); err != nil {
		log.Println(ctx, "Unmarshal err: %v. Raw response: %s", err, response.Content[0].Text)
		return nil, err
	}

	return &translatedChapter, nil
}
