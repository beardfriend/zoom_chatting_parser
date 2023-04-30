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
		message := `Reacted to "그래서 저는 다른사람이 분석한거 먼저..." with 👏🏻`
		emoji, sentence := ReactionExtractFunc(parser, message)
		assert.Equal(t, "👏🏻", emoji)
		assert.Equal(t, "그래서 저는 다른사람이 분석한거 먼저", sentence)
	})
}
