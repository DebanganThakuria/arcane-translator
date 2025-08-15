package llm

import (
	"context"
	"encoding/json"
	"log"

	"backend/models"

	"google.golang.org/genai"
)

const (
	ResponseMimeType     = "application/json"
	GenerateContentModel = "gemini-2.5-flash"
)

type geminiClientImpl struct {
	geminiClient *genai.Client
}

func (g geminiClientImpl) TranslateNovelDetails(ctx context.Context, webpageContent string) (*models.NovelDetails, error) {
	prompt := `
	You are a professional translator for webnovels.
	Please extract and translate information from this novel page content.

	Page content: ` + webpageContent + `

	Return ONLY a valid JSON object with this exact structure:
	{
		"novel_title_original": "Original title in the source language",
		"novel_title_translated": "Translated title in English",
		"novel_summary_translated": "Translated summary in English in HTML format with paragraph tags. Please ensure that the summary has valid HTML tags for rendering on the frontend.",
		"possible_novel_genres": ["Genre1", "Genre2", ...],
		"number_of_chapters": "total number of chapters in original as integer",
		"status": "Ongoing" or "Completed" or "Unknown"
	}

	Do not include any commentary, explanation, or preamble. Only return the JSON object.`

	response, err := g.geminiClient.Models.GenerateContent(ctx,
		GenerateContentModel,
		[]*genai.Content{genai.NewContentFromText(prompt, genai.RoleUser)},
		&genai.GenerateContentConfig{
			Temperature:      genai.Ptr(float32(0.3)),
			MaxOutputTokens:  65000,
			ResponseMIMEType: ResponseMimeType,
			ResponseSchema: &genai.Schema{
				Type: "object",
				Properties: map[string]*genai.Schema{
					"novel_title_original": {
						Type:        "string",
						Description: "Original title in the source language",
						Nullable:    genai.Ptr(false),
					},
					"novel_title_translated": {
						Type:        "string",
						Description: "Translated title in English",
						Nullable:    genai.Ptr(false),
					},
					"novel_summary_translated": {
						Type:        "string",
						Description: "Translated summary in English in HTML format with paragraph tags. Please ensure that the summary has valid HTML tags for rendering on the frontend.",
						Nullable:    genai.Ptr(false),
					},
					"possible_novel_genres": {
						Type:        "array",
						Description: "List of possible genres for the novel",
						Items: &genai.Schema{
							Type: "string",
						},
						Nullable: genai.Ptr(false),
					},
					"number_of_chapters": {
						Type:        "integer",
						Description: "Total number of chapters in original as integer",
						Nullable:    genai.Ptr(false),
					},
					"status": {
						Type:        "string",
						Description: "Status of the novel, can be 'Ongoing', 'Completed', or 'Unknown'",
						Nullable:    genai.Ptr(false),
					},
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	log.Println(response.Text())

	var novelDetails models.NovelDetails
	if err = json.Unmarshal([]byte(response.Text()), &novelDetails); err != nil {
		return nil, err
	}

	return &novelDetails, nil
}

func (g geminiClientImpl) TranslateNovelChapter(ctx context.Context, novelGenres []string, webpageContent string) (*models.TranslatedChapter, error) {
	prompt := `
        You are the best webnovel translator and editor, capable of producing the highest quality work.
		Your task is translating and polishing the following Webnovel chapter into flawless English, ensuring perfect grammar and language. Translate all original language object names, including places, abilities, techniques, and other cultural references, into English.
		Your goal is to craft an engaging English version while preserving the original content. The author of the original text may have limited writing skills, so take the liberty to polish the content in your translation wherever you judge necessary. Pay special attention to the dialogues, ensuring they flow smoothly and sound lifelike. Do not add anything extra anywhere from your end.
		Finally, you mustn't lose any content from the original during the translation process. 
		I trust you to provide the best possible results. Please translate the full chapter as per these guidelines.

        Currently known novel genres: ` + models.GenresToString(novelGenres) + `
		Chapter Content in the source language: ` + webpageContent + `

        Return ONLY a valid JSON object with this exact structure:
        {{
            "translated_chapter_title": "The translated title of the chapter. This is important. Never miss it.",
			"original_chapter_title": "The original title in the source language.",
            "translated_chapter_contents": "The translated content of the full chapter in HTML format with paragraph tags. Please ensure that the chapter content has valid HTML tags for rendering on the frontend. And, most importantly, ensure that the full chapter content is included in the response.",
            "possible_new_genres": "Array of any genres you detect that aren't already known (leave empty if none). Response should be like this ["Genre1", "Genre2", ...]"
        }}
	
		Additional instructions. Please follow these instructions carefully:
		- Please terminate the translation and return when the chapter ends.
		- Do not miss any of the fields. This is very important. 
		- Do not include any commentary, explanation, or preamble. Only return the JSON object.
	`

	response, err := g.geminiClient.Models.GenerateContent(ctx,
		GenerateContentModel,
		[]*genai.Content{genai.NewContentFromText(prompt, genai.RoleUser)},
		&genai.GenerateContentConfig{
			Temperature:      genai.Ptr(float32(0.3)),
			MaxOutputTokens:  65000,
			ResponseMIMEType: ResponseMimeType,
			ResponseSchema: &genai.Schema{
				Type: "object",
				Properties: map[string]*genai.Schema{
					"translated_chapter_title": {
						Type:        "string",
						Description: "The translated title of the chapter.",
						Nullable:    genai.Ptr(false),
					},
					"original_chapter_title": {
						Type:        "string",
						Description: "The original title in the source language.",
						Nullable:    genai.Ptr(false),
					},
					"translated_chapter_contents": {
						Type:        "string",
						Description: "The translated content of the chapter in HTML format with paragraph tags. Please ensure that the chapter content has valid HTML tags for rendering on the frontend.",
						Nullable:    genai.Ptr(false),
					},
					"possible_new_genres": {
						Type:        "array",
						Description: "Array of any genres you detect that aren't already known (leave empty if none). Response should be like this [\"Genre1\", \"Genre2\", ...]",
						Items: &genai.Schema{
							Type: "string",
						},
						Nullable: genai.Ptr(false),
					},
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	log.Println(response.Text())

	var translatedChapter models.TranslatedChapter
	if err = json.Unmarshal([]byte(response.Text()), &translatedChapter); err != nil {
		return nil, err
	}

	return &translatedChapter, nil
}
