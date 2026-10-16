package llm

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

// fakeConverser records the request and returns a canned reply.
type fakeConverser struct {
	got *bedrockruntime.ConverseInput
	out *bedrockruntime.ConverseOutput
	err error
}

func (f *fakeConverser) Converse(_ context.Context, in *bedrockruntime.ConverseInput, _ ...func(*bedrockruntime.Options)) (*bedrockruntime.ConverseOutput, error) {
	f.got = in
	return f.out, f.err
}

func reply(stop types.StopReason, blocks ...types.ContentBlock) *bedrockruntime.ConverseOutput {
	return &bedrockruntime.ConverseOutput{
		StopReason: stop,
		Output: &types.ConverseOutputMemberMessage{Value: types.Message{
			Role:    types.ConversationRoleAssistant,
			Content: blocks,
		}},
	}
}

const chapterJSON = `{"translated_chapter_title":"Chapter 70","original_chapter_title":"第70章","translated_chapter_contents":"<p>Text</p>","possible_new_genres":[]}`

func TestBedrockTranslateNovelChapter(t *testing.T) {
	fake := &fakeConverser{out: reply(types.StopReasonEndTurn,
		&types.ContentBlockMemberReasoningContent{Value: &types.ReasoningContentBlockMemberReasoningText{
			Value: types.ReasoningTextBlock{Text: aws.String("thinking about it")},
		}},
		&types.ContentBlockMemberText{Value: "```json\n" + chapterJSON + "\n```"},
	)}
	client := bedrockClientImpl{client: fake, modelID: "zai.glm-4.7"}

	got, err := client.TranslateNovelChapter(context.Background(), []string{"Fantasy"}, "<html>page</html>")
	if err != nil {
		t.Fatalf("TranslateNovelChapter() error = %v", err)
	}
	if got.TranslatedChapterTitle != "Chapter 70" || got.TranslatedChapterContents != "<p>Text</p>" {
		t.Errorf("TranslateNovelChapter() = %+v, want the fields from the reply", got)
	}

	if id := aws.ToString(fake.got.ModelId); id != "zai.glm-4.7" {
		t.Errorf("ModelId = %q, want zai.glm-4.7", id)
	}
	system := fake.got.System[0].(*types.SystemContentBlockMemberText).Value
	if system != chapterPrompt([]string{"Fantasy"}) {
		t.Error("system prompt is not chapterPrompt(genres)")
	}
	user := fake.got.Messages[0].Content[0].(*types.ContentBlockMemberText).Value
	if user != "<html>page</html>" {
		t.Errorf("user message = %q, want the page HTML", user)
	}
}

func TestBedrockConverseErrors(t *testing.T) {
	tests := []struct {
		name    string
		fake    *fakeConverser
		wantErr string
	}{
		{"api error", &fakeConverser{err: errors.New("ResourceNotFoundException")}, "ResourceNotFoundException"},
		{"truncated reply", &fakeConverser{out: reply(types.StopReasonMaxTokens, &types.ContentBlockMemberText{Value: "{"})}, "cut off"},
		{"no text", &fakeConverser{out: reply(types.StopReasonEndTurn)}, "no text"},
		{"unparseable json", &fakeConverser{out: reply(types.StopReasonEndTurn, &types.ContentBlockMemberText{Value: `{"translated_chapter_title": \p}`})}, "unmarshal"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := bedrockClientImpl{client: tt.fake, modelID: "zai.glm-4.7"}
			_, err := client.TranslateNovelChapter(context.Background(), nil, "page")
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("TranslateNovelChapter() error = %v, want one containing %q", err, tt.wantErr)
			}
		})
	}
}
