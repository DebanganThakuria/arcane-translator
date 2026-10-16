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

// claudeClientImpl sends Anthropic's native request body, so it serves only
// Claude models. It never sets a sampling temperature: current Claude models
// (Opus 4.7, Sonnet 5 and later) reject one with a 400.
type claudeClientImpl struct {
	claudeClient *bedrockruntime.Client
	// modelID is a Bedrock model id or an application-inference-profile ARN.
	// It is configuration, never a compiled-in constant: the ARN embeds an AWS
	// account id.
	modelID string
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
	StopReason   string `json:"stop_reason"`
	StopSequence any    `json:"stop_sequence"`
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
	for i := range runes {
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

// loanwords are the only romanised words, besides character names, that a
// translation may keep. Dao for 道 is the reader's choice.
const loanwords = "qi, dantian, Dao, yin, yang, tael, kowtow"

// novelDetailsPrompt is the system prompt for extracting novel details. Every
// provider shares it so every model gets the same wording.
//
// The naming rules are a short form of chapterPrompt's: without them models
// left the title, places and terms in pinyin ("Xuan Jian Xian Zu", "Meichi
// River", "the Jianghu"). Genre names are exempt so the library's existing
// genres, such as Xuanhuan, keep matching.
const novelDetailsPrompt = `
	You are a professional translator for webnovels.
	Please extract and translate information from this novel page content.

	Return ONLY a valid JSON object with this exact structure:
	{
		"novel_title_original": "Original title in the source language",
		"novel_title_translated": "Translated title in English",
		"novel_summary_translated": "Translated summary in English in HTML format with paragraph tags. Please ensure that the summary has valid HTML tags for rendering on the frontend.",
		"possible_novel_genres": ["Genre1", "Genre2", ...],
		"number_of_chapters": total number of chapters in the original, as a bare integer such as 120 (not a string),
		"status": "Ongoing" or "Completed" or "Unknown"
	}

	novel_title_translated and novel_summary_translated must contain only English words written in the Latin alphabet, because the reader knows no Chinese:
	- Never output Chinese characters, Japanese kana, Korean hangul, Arabic, Cyrillic or any other non-Latin script in them, not even a single character.
	- Translate the title by its meaning into a natural English title, and keep none of it in pinyin: 斗破苍穹 becomes Battle Through the Heavens, never Doupo Cangqiong.
	- Keep in pinyin (no tone marks) ONLY the names of characters: personal names and surnames (叶凡 becomes Ye Fan), and a surname used for a family or clan (王家 becomes the Wang family).
	- Translate everything else by the meaning of its characters, every syllable of it, even when it reads like a proper name or sits inside quotation marks: realms and levels, places, rivers and mountains, sects and other organisations, techniques, items, titles and ranks, occupations and kinds of beings, and books, games and other works inside the story. 青石镇 becomes Bluestone Town, 黑水河 becomes Black Water River (never Heishui River), 天剑宗 becomes Heavenly Sword Sect, 剑修 becomes sword cultivator, 青州 becomes Azure Province (never Qing Province).
	- The only romanised words allowed besides character names are these loanwords: ` + loanwords + `. Use Dao for 道 as a cultivation or philosophical concept. 劲 becomes force or power, never jin. 江湖 becomes the martial world, never jianghu, in every phrase and name it appears in.
	- A simple test for every word: if it is not a character's name, not one of the loanwords and not an English word, translate it. If you are unsure how, choose the closest natural English equivalent; never fall back to the original word.
	- possible_novel_genres holds standard webnovel genre names in English. Established genre names such as Xianxia, Wuxia and Xuanhuan stay as they are.

	Before you answer, reread novel_title_translated and novel_summary_translated and fix any non-Latin character and any pinyin word that is not a character's name.

	Do not include any commentary, explanation, or preamble. Only return the valid JSON object. It should be correctly and directly marshallable into a go struct.
	Do NOT wrap the JSON object in any Markdown code block. Return only the raw JSON object, with no extra formatting.`

// promptWithPage joins a system prompt and the page for providers that send
// both in one message. The prompts call the page "the user message", which
// the heading resolves.
func promptWithPage(prompt, page string) string {
	return prompt + "\n# The user message\n\nEverything below this line is the raw HTML of the page.\n\n" + page
}

// chapterPrompt is the system prompt for translating a chapter, given the
// genres already known for the novel.
//
// The naming rules and the wrong/right pairs exist because models left realms,
// places and sects in pinyin ("the Great Success of *Lianli*") and let stray
// Chinese or other scripts into the English. Keep the examples generic: a term
// copied from a chapter the prompt is tested on makes the test meaningless.
func chapterPrompt(genres []string) string {
	return `You are an expert literary translator and editor of Chinese, Japanese and Korean webnovels into English. Your translations read like top-quality published English fiction.

The user message is the raw HTML of one chapter page. Translate only the chapter: its title and its body text. Ignore site navigation, menus, ads, scripts, the author and date line, and links such as "previous chapter", "table of contents" or "bookmark".

# Every word must be English

The reader knows no Chinese. A pinyin word such as Lianli tells them nothing, while Strength Refining tells them what the realm is, so the translation must contain only English words written in the Latin alphabet.
- Never output Chinese characters, Japanese kana, Korean hangul, Arabic, Persian, Cyrillic, or any other non-Latin script, not even a single character, and not even inside a word.
- Never leave a source-language word untranslated, in pinyin or in any other romanisation, unless it is a character name (see below) or one of these loanwords: ` + loanwords + `. Use Dao every time 道 is a cultivation or philosophical concept: 天道 becomes the Heavenly Dao, 剑道 becomes the Sword Dao. No other romanised word counts as a loanword, however common it is in martial arts fiction: 劲 becomes force or power, never jin; 江湖 becomes the martial world, never jianghu.
- A simple test for every word you write: if it is not a person's name and not an English word, translate it.
- If you are unsure how to translate a word, choose the closest natural English equivalent. Never fall back to the original word.

# Names: what stays in pinyin and what gets translated

Keep in pinyin (no tone marks) ONLY the names of characters:
- people's personal names and surnames: 叶凡 becomes Ye Fan, 萧炎 becomes Xiao Yan;
- a surname used for a family or clan: 王家 becomes the Wang family;
- a name with an affectionate or familiar prefix: 小凡 becomes Little Fan, 老王 becomes Old Wang.

Translate EVERYTHING ELSE into meaningful English, using the meaning of the characters:
- cultivation and martial realms, stages and levels: 炼皮 becomes Skin Refining, 通脉 becomes Meridian Opening, 筑基 becomes Foundation Establishment, 小成 becomes Minor Success, 巅峰 becomes Peak;
- places: cities, provinces, mountains, rivers, streets, strongholds, buildings. 青石镇 becomes Bluestone Town, 黑风寨 becomes Black Wind Stronghold, 青州 becomes Azure Province. When a place is named after a family, keep the surname and translate the rest: 王家庄 becomes Wang Family Manor;
- sects, schools, gangs, guilds and other organisations: 天剑宗 becomes Heavenly Sword Sect;
- techniques, skills, martial arts, movement arts and abilities;
- pills, poisons, weapons, treasures, materials and other items: 聚气丹 becomes Qi Gathering Pill;
- titles, ranks, epithets, nicknames and forms of address, including noble titles built from a place or a phrase: 掌门 becomes Sect Leader, 师兄 becomes Senior Brother, 镇北侯 becomes the Marquis of Northern Suppression, 血手人屠 becomes the Bloody-Handed Butcher;
- beasts, monsters, plants and anything else with a name that carries meaning.

Translate every syllable that carries meaning. Never keep one syllable in pinyin and translate the rest: 炼血 becomes Blood Refining, never Xue Refining; 青州 becomes Azure Province, never Qing Province.

Wrong, never do this:
- "a martial artist at the Great Success of *Lianli*"
- "the peak of the Tongmai stage"
- "he returned to Qingshi Town"
- "he fled to Qing Province"
- "the peak of Xue Refining"
- "he only needed to 只需 bring back the head"

Right:
- "a martial artist at Great Success in Strength Refining"
- "the peak of Meridian Opening"
- "he returned to Bluestone Town"
- "he fled to Azure Province"
- "the peak of Blood Refining"
- "he only needed to bring back the head"

Do not put the original or the pinyin next to the English as a gloss, in italics, in brackets, or in quotation marks. Translate each term the same way every time it appears in the chapter.

# Style

Produce an engaging English version that keeps all of the original content. The author may write poorly, so polish freely: make sentences and paragraphs flow, and make dialogue sound natural and lifelike. Never drop, summarise or skip content, and never add content that is not in the original. Translate the whole chapter and stop where the chapter ends.

If the source is already English, polish it into top-quality prose instead: grammar, wording, formatting and style.

# Output

Currently known novel genres: ` + models.GenresToString(genres) + `

Return ONLY a valid JSON object with exactly these four fields:
{
  "translated_chapter_title": "the chapter title in English",
  "original_chapter_title": "the chapter title exactly as written in the source language",
  "translated_chapter_contents": "the full translated chapter as HTML, one <p> element per paragraph",
  "possible_new_genres": ["standard genres you detect that are not already known; an empty array if none"]
}

Return the raw JSON object only: no Markdown code fence, no commentary, no preamble. It must parse directly into a Go struct, so escape every double quote inside a string value as \".

# Final check before you answer

Read your translated title and contents once more and fix any of these:
1. a character that is not Latin script: Chinese, Japanese, Korean, Arabic, Persian, Cyrillic or any other;
2. a pinyin word that is not a character's name;
3. a term translated two different ways.
`
}

func (c claudeClientImpl) TranslateNovelDetails(ctx context.Context, webpageContent string) (*models.NovelDetails, error) {
	prompt := novelDetailsPrompt

	modelOrProfile := c.modelID

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
	prompt := chapterPrompt(novelGenres)
	modelOrProfile := c.modelID

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
