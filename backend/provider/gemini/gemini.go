package gemini

import (
	"context"
	"encoding/json"
	"os"

	"backend/models"

	"google.golang.org/genai"
)

const (
	ResponseMimeType     = "application/json"
	GenerateContentModel = "gemini-2.5-flash-preview-05-20"
)

type IClient interface {
	TranslateNovelDetails(ctx context.Context, webpageContent string) (*models.NovelDetails, error)
	TranslateNovelChapter(ctx context.Context, novelGenres []string, webpageContent string) (*models.TranslatedChapter, error)
}

var geminiClient IClient

type geminiClientImpl struct {
	geminiClient *genai.Client
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
}

func GetClient() IClient {
	return geminiClient
}

func (g geminiClientImpl) TranslateNovelDetails(ctx context.Context, webpageContent string) (*models.NovelDetails, error) {
	prompt := `You are a professional translator for webnovels.
	Please extract and translate information from this novel page content.

	Page content: ` + webpageContent + `

	Return ONLY a valid JSON object with this exact structure:
	{
		"novel_title_original": "Original title in the source language",
		"novel_title_translated": "Translated title in English",
		"novel_summary_translated": "Translated summary in English. You may use html paragraph, bold or italic tags for formatting as per the original text.",
		"possible_novel_genres": ["Genre1", "Genre2", ...],
		"number_of_chapters": "total number of chapters in original as integer",
		"status": "Ongoing" or "Completed" or "Unknown"
	}

	Do not include any commentary, explanation, or preamble. Only return the JSON object.`

	response, err := g.geminiClient.Models.GenerateContent(ctx,
		GenerateContentModel,
		[]*genai.Content{genai.NewContentFromText(prompt, genai.RoleModel)},
		&genai.GenerateContentConfig{
			ResponseMIMEType: ResponseMimeType,
			ResponseSchema: &genai.Schema{
				Type: "object",
				Properties: map[string]*genai.Schema{
					"novel_title_original": {
						Type:        "string",
						Description: "Original title in the source language",
					},
					"novel_title_translated": {
						Type:        "string",
						Description: "Translated title in English",
					},
					"novel_summary_translated": {
						Type:        "string",
						Description: "Translated summary in English. You may use html paragraph, bold or italic tags for formatting as per the original text.",
					},
					"possible_novel_genres": {
						Type:        "array",
						Description: "List of possible genres for the novel",
						Items: &genai.Schema{
							Type: "string",
						},
					},
					"number_of_chapters": {
						Type:        "integer",
						Description: "Total number of chapters in original as integer",
					},
					"status": {
						Type:        "string",
						Description: "Status of the novel, can be 'Ongoing', 'Completed', or 'Unknown'",
					},
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	var novelDetails models.NovelDetails
	if err = json.Unmarshal([]byte(response.Text()), &novelDetails); err != nil {
		return nil, err
	}

	return &novelDetails, nil
}

func (g geminiClientImpl) TranslateNovelChapter(ctx context.Context, novelGenres []string, webpageContent string) (*models.TranslatedChapter, error) {
	prompt := `
		TASK:
        You are the best webnovel translator and editor, capable of producing the highest quality work. Your task is translating and polishing the following Webnovel chapter into flawless English, ensuring perfect grammar and language. Translate all original language object names, including places, abilities, techniques, and other cultural references, into English. Provide necessary context and clarification for Western readers where needed.
        Your goal is to craft an engaging, culturally sensitive English version while preserving the original storytelling style. The author of the original text may have limited writing skills, so take the liberty to polish the content in your translation. Pay special attention to the dialogues, ensuring they flow smoothly and sound lifelike. Think of yourself as the author.
        You mustn't lose any plot points or system messages during the translation process. Correct any logical errors to ensure the story makes sense to the reader. I trust you to provide the best possible results.

        NOVEL CONTEXT:
        Genre: ` + models.GenresToString(novelGenres) + `

        Return ONLY a valid JSON object with this exact structure:
        {{
            "translated_chapter_title": "The translated title",
			"original_chapter_title": "The original title in the source language",
            "translated_chapter_contents": "The translated content in HTML format with paragraph tags",
            "possible_new_genres": "Array of any genres you detect that aren't already listed (leave empty if none). Response should be like this ["Genre1", "Genre2", ...]"
        }}

        Content: ` + webpageContent

	response, err := g.geminiClient.Models.GenerateContent(ctx,
		GenerateContentModel,
		[]*genai.Content{genai.NewContentFromText(prompt, genai.RoleModel)},
		&genai.GenerateContentConfig{
			ResponseMIMEType: ResponseMimeType,
			ResponseSchema: &genai.Schema{
				Type: "object",
				Properties: map[string]*genai.Schema{
					"translated_chapter_title": {
						Type:        "string",
						Description: "Translated chapter title in English",
					},
					"original_chapter_title": {
						Type:        "string",
						Description: "Original chapter title in the source language",
					},
					"translated_chapter_contents": {
						Type:        "string",
						Description: "Translated chapter contents in English. You may use html paragraph, bold or italic tags for formatting as per the original text.",
					},
					"possible_new_genres": {
						Type:        "array",
						Description: "List of possible new genres for the chapter",
						Items: &genai.Schema{
							Type: "string",
						},
					},
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	var translatedChapter models.TranslatedChapter
	if err = json.Unmarshal([]byte(response.Text()), &translatedChapter); err != nil {
		return nil, err
	}

	return &translatedChapter, nil
}
