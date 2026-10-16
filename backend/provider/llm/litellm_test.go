package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/sashabaranov/go-openai"
)

// litellmServer serves chat completions. reject, when set, answers a request
// with that status and body instead.
type litellmServer struct {
	bodies []map[string]any
	reject func(body map[string]any) (int, string)
	reply  string
	finish openai.FinishReason
}

func (s *litellmServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body map[string]any
	_ = json.NewDecoder(r.Body).Decode(&body)
	s.bodies = append(s.bodies, body)

	w.Header().Set("Content-Type", "application/json")
	if s.reject != nil {
		if status, msg := s.reject(body); status != 0 {
			w.WriteHeader(status)
			_, _ = w.Write([]byte(msg))
			return
		}
	}
	_ = json.NewEncoder(w).Encode(openai.ChatCompletionResponse{
		Choices: []openai.ChatCompletionChoice{{Message: openai.ChatCompletionMessage{Role: "assistant", Content: s.reply}, FinishReason: s.finish}},
	})
}

func newLiteLLMTestClient(t *testing.T, s *litellmServer) litellmClientImpl {
	t.Helper()
	return newLiteLLMTestClientFor(t, s, "glm-5")
}

func newLiteLLMTestClientFor(t *testing.T, s *litellmServer, model string) litellmClientImpl {
	t.Helper()
	srv := httptest.NewServer(s)
	t.Cleanup(srv.Close)

	cfg := openai.DefaultConfig("test-key")
	cfg.BaseURL = srv.URL
	return *newLiteLLMClient(cfg, model, 0.2)
}

func TestLiteLLMCheapMode(t *testing.T) {
	off := map[string]any{"type": "between_tools"}
	medium := map[string]any{"effort": "medium"}
	tests := []struct {
		model              string
		thinking, settings any
	}{
		{"claude-sonnet-5-5", off, medium},
		{"anthropic/claude-sonnet-5.5", off, medium},
		// Every other model rejects between_tools.
		{"claude-sonnet-5", nil, nil},
		{"claude-opus-5-5", nil, nil},
		{"glm-5", nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			s := &litellmServer{reply: chapterJSON}
			client := newLiteLLMTestClientFor(t, s, tt.model)

			if _, err := client.TranslateNovelChapter(context.Background(), []string{"Fantasy"}, "page"); err != nil {
				t.Fatalf("TranslateNovelChapter() error = %v", err)
			}
			body := s.bodies[0]
			if got := body["thinking"]; !reflect.DeepEqual(got, tt.thinking) {
				t.Errorf("thinking = %v, want %v", got, tt.thinking)
			}
			if got := body["output_config"]; !reflect.DeepEqual(got, tt.settings) {
				t.Errorf("output_config = %v, want %v", got, tt.settings)
			}
			if body["model"] != tt.model || body["temperature"] == nil || len(body["messages"].([]any)) != 1 {
				t.Errorf("the rest of the request changed: %v", body)
			}
		})
	}
}

func TestLiteLLMTranslateNovelChapter(t *testing.T) {
	s := &litellmServer{reply: "```json\n" + chapterJSON + "\n```"}
	client := newLiteLLMTestClient(t, s)

	got, err := client.TranslateNovelChapter(context.Background(), []string{"Fantasy"}, "<html>page</html>")
	if err != nil {
		t.Fatalf("TranslateNovelChapter() error = %v", err)
	}
	if got.TranslatedChapterTitle != "Chapter 70" {
		t.Errorf("TranslateNovelChapter() = %+v, want the fields from the reply", got)
	}

	body := s.bodies[0]
	if temp, _ := body["temperature"].(float64); temp < 0.19 || temp > 0.21 {
		t.Errorf("temperature = %v, want 0.2", body["temperature"])
	}
	if got, _ := body["max_completion_tokens"].(float64); got != litellmMaxTokens {
		t.Errorf("max_completion_tokens = %v, want %d", body["max_completion_tokens"], litellmMaxTokens)
	}
	messages, _ := body["messages"].([]any)
	if len(messages) != 1 {
		t.Fatalf("got %d messages, want the prompt and page in one user message", len(messages))
	}
	msg, _ := messages[0].(map[string]any)
	if msg["role"] != "user" || msg["content"] != promptWithPage(chapterPrompt([]string{"Fantasy"}), "<html>page</html>") {
		t.Error("user message is not chapterPrompt(genres) followed by the page")
	}
}

func TestLiteLLMTranslateNovelDetailsUsesSharedPrompt(t *testing.T) {
	s := &litellmServer{reply: `{"novel_title_original":"书","novel_title_translated":"Book","novel_summary_translated":"<p>S</p>","number_of_chapters":3,"status":"Ongoing"}`}
	client := newLiteLLMTestClient(t, s)

	got, err := client.TranslateNovelDetails(context.Background(), "<html>index</html>")
	if err != nil {
		t.Fatalf("TranslateNovelDetails() error = %v", err)
	}
	if got.NovelTitleTranslated != "Book" || got.NumberOfChapters != 3 {
		t.Errorf("TranslateNovelDetails() = %+v, want the fields from the reply", got)
	}

	msg := s.bodies[0]["messages"].([]any)[0].(map[string]any)
	if msg["content"] != promptWithPage(novelDetailsPrompt, "<html>index</html>") {
		t.Error("user message is not novelDetailsPrompt followed by the page")
	}
}

// OpenAI's reasoning models answer a non-default temperature with this error.
const temperatureRejected = `{"error":{"message":"Unsupported value: 'temperature' does not support 0.2 with this model. Only the default (1) value is supported.","type":"invalid_request_error","param":"temperature","code":"unsupported_value"}}`

func TestLiteLLMRetriesWithoutRejectedTemperature(t *testing.T) {
	s := &litellmServer{
		reply: chapterJSON,
		reject: func(body map[string]any) (int, string) {
			if _, ok := body["temperature"]; ok {
				return http.StatusBadRequest, temperatureRejected
			}
			return 0, ""
		},
	}
	client := newLiteLLMTestClient(t, s)

	if _, err := client.TranslateNovelChapter(context.Background(), nil, "page"); err != nil {
		t.Fatalf("TranslateNovelChapter() error = %v", err)
	}
	if len(s.bodies) != 2 {
		t.Fatalf("got %d requests, want the rejected one and a retry", len(s.bodies))
	}
	if _, ok := s.bodies[1]["temperature"]; ok {
		t.Error("retry still sends temperature")
	}
}

// go-openai itself refuses a non-default temperature for gpt-5 model names,
// before any request is sent.
func TestLiteLLMDropsTemperatureForReasoningModelNames(t *testing.T) {
	s := &litellmServer{reply: chapterJSON}
	client := newLiteLLMTestClient(t, s)
	client.model = "gpt-5"

	if _, err := client.TranslateNovelChapter(context.Background(), nil, "page"); err != nil {
		t.Fatalf("TranslateNovelChapter() error = %v", err)
	}
	if len(s.bodies) != 1 {
		t.Fatalf("got %d requests, want 1", len(s.bodies))
	}
	if _, ok := s.bodies[0]["temperature"]; ok {
		t.Error("request for gpt-5 still sends temperature")
	}
}

func TestLiteLLMDoesNotRetryOtherErrors(t *testing.T) {
	s := &litellmServer{
		reject: func(map[string]any) (int, string) {
			return http.StatusBadRequest, `{"error":{"message":"bad model","type":"invalid_request_error","param":"model","code":null}}`
		},
	}
	client := newLiteLLMTestClient(t, s)

	_, err := client.TranslateNovelChapter(context.Background(), nil, "page")
	if err == nil || !strings.Contains(err.Error(), "bad model") {
		t.Errorf("TranslateNovelChapter() error = %v, want the API error", err)
	}
	if len(s.bodies) != 1 {
		t.Errorf("got %d requests, want 1", len(s.bodies))
	}
}

// A reply cut off at the token limit ends mid-JSON; say so rather than
// reporting a parse error.
func TestLiteLLMTruncatedReply(t *testing.T) {
	s := &litellmServer{reply: `{"translated_chapter_title":"Chapter 93","translated_chapter_contents":"<p>Half`, finish: openai.FinishReasonLength}
	client := newLiteLLMTestClient(t, s)

	_, err := client.TranslateNovelChapter(context.Background(), nil, "page")
	if err == nil || !strings.Contains(err.Error(), "cut off") {
		t.Errorf("TranslateNovelChapter() error = %v, want one saying the reply was cut off", err)
	}
}
