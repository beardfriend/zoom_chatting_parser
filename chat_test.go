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

func TestModule(t *testing.T) {
	extractReaction()
}
