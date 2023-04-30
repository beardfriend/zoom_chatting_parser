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
