package chat_test

import (
	"os"
	"testing"

	z "github.com/beardfriend/zoom_chat_parser"

	"github.com/stretchr/testify/assert"
)

func TestFileNotExist(t *testing.T) {
	parser := z.NewParser()

	t.Run("nil file", func(t *testing.T) {
		err, _ := parser.Parse(nil)
		assert.NotEmpty(t, err)
	})

	t.Run("no such file", func(t *testing.T) {
		// load nonexistend File
		file, _ := os.Open("nonexistent_file.txt")

		// Parse
		err, _ := parser.Parse(file)

		assert.NotEmpty(t, err)
	})
}

func TestFilExist(t *testing.T) {
	parser := z.NewParser()
	t.Run("demo", func(t *testing.T) {
		file, _ := os.Open("assets/test.txt")
		err, _ := parser.Parse(file)
		if err != nil {
			assert.NoError(t, err)
		}
	})
}

func TestExtractReaction(t *testing.T) {
	parser := z.NewParser()
	ReactionExtractFunc := z.ExportExtractReaction

	t.Run("normal message", func(t *testing.T) {
		message := `	Reacted to "좋은 아침이예요!" with 🙌`
		emoji, sentence := ReactionExtractFunc(parser, message)
		assert.Equal(t, "🙌", emoji)
		assert.Equal(t, "좋은 아침이예요!", sentence)
	})

	t.Run("over char message", func(t *testing.T) {
		message := `	Reacted to "그래서 저는 다른사람이 분석한거 먼저..." with 👏🏻`
		emoji, sentence := ReactionExtractFunc(parser, message)
		assert.Equal(t, "👏🏻", emoji)
		assert.Equal(t, "그래서 저는 다른사람이 분석한거 먼저", sentence)
	})
}

func TestExtractReply(t *testing.T) {
	parser := z.NewParser()
	ReplyExtractFunc := z.ExportExtractReply

	t.Run("url", func(t *testing.T) {
		message := `	Replying to "https://plotly.com/p..."`
		sentence := ReplyExtractFunc(parser, message)

		assert.Equal(t, "https://plotly.com/p", sentence)
	})

	t.Run("normal text over char", func(t *testing.T) {
		message := `	Replying to "파일이 다운로드 안되네요. 다운되자마..."`
		sentence := ReplyExtractFunc(parser, message)
		assert.Equal(t, "파일이 다운로드 안되네요. 다운되자마", sentence)
	})
}

func TestExtractRemove(t *testing.T) {
	parser := z.NewParser()
	RemoveExtractFunc := z.ExportExtractRemove

	t.Run("normal text char over", func(t *testing.T) {
		message := `	Removed a 👍 reaction from "plt.rcParams['font.f..."`
		emoji, sentence := RemoveExtractFunc(parser, message)
		assert.Equal(t, "👍", emoji)
		assert.Equal(t, "plt.rcParams['font.f", sentence)
	})
}

func TestFindIdFromChatHistory(t *testing.T) {
	parser := z.NewParser()
	result := z.Result{
		ZoomChatHistory: []*z.ZoomChatHistory{
			{
				Id:           1,
				TextType:     z.NormalText,
				SenderName:   "박세훈",
				ReceiverName: "모두에게",
				Text:         "좋은 아침이예요!",
				ChatedAt:     "09:00:01",
			},
			{
				Id:           2,
				TextType:     z.NormalText,
				SenderName:   "김또깡",
				ReceiverName: "모두에게",
				Text:         "주피터돌리다가 크롬이 구역질해요",
				ChatedAt:     "09:00:03",
			},
			{
				Id:           3,
				TextType:     z.NormalText,
				SenderName:   "아라이",
				ReceiverName: "모두에게",
				Text:         "예압",
				ChatedAt:     "09:00:01",
			},
			{
				Id:           4,
				TextType:     z.NormalText,
				SenderName:   "오우상",
				ReceiverName: "모두에게",
				Text:         "겨울인데 여름처럼 더워지고",
				ChatedAt:     "09:00:01",
			},
		},
	}
	t.Run("텍스트가 왼전 일치하는 경우", func(t *testing.T) {
		FindIdFromChatHistory := z.ExportFindIdFromChatHistoryByText(parser, "좋은 아침이예요!", uint(5), result.ZoomChatHistory)
		assert.Equal(t, uint(1), *FindIdFromChatHistory)
	})

	t.Run("텍스트 일부가 포함되어 있는 경우", func(t *testing.T) {
		// Todo: 찾으려는 텍스트 이전에
		FindIdFromChatHistory := z.ExportFindIdFromChatHistoryByText(parser, "좋은", uint(5), result.ZoomChatHistory)
		assert.Equal(t, uint(1), *FindIdFromChatHistory)
	})
}
