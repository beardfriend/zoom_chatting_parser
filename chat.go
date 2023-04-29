package chat

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type ChatContentType uint

const (
	NormalText ChatContentType = 1 + iota
	Reaction
	Reply
	Remove
)

type ZoomChatHistory struct {
	Id uint

	ReactId *uint

	ReplyId *uint

	TextType ChatContentType

	SenderName string

	ReceiverName string

	Text string

	ReactCount *uint

	ReplyCount *uint

	// format HH:MM:ss
	ChatedAt string
}

type Result struct {
	ZoomChatHistory []*ZoomChatHistory
	Statistic       struct {
		MissingReactionIds []int
		MissingReplyIds    []int
	}
}

var ErrorNoFile = errors.New("no such file")

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(file *os.File) (err error, result Result) {
	if file == nil {
		err = ErrorNoFile
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var index, textCursor, searchCursor uint

	for scanner.Scan() {
		line := scanner.Text()

		// MetaData
		if !strings.Contains(line, "\t") {
			splited := strings.Split(line, " ")

			result.ZoomChatHistory = append(result.ZoomChatHistory, &ZoomChatHistory{
				Id:           index,
				ChatedAt:     splited[0],
				SenderName:   splited[2],
				ReceiverName: splited[4],
			})

			index++
			textCursor = 0
		} else {

			text := strings.TrimPrefix(line, "\t")
			category := categorizeMessage(text)

			// 리엑션일 때는 리엑션이 어느 글에서 했는지 판단
			if category == Reaction {
				emoji, sentence := extractReaction(text)

				searchCursor = index - 1

				// 어느 글에서 리엑션을 했는지 판단하는
				for searchCursor > 0 && searchCursor > index-20 {

					beforeChat := result.ZoomChatHistory[searchCursor]

					if strings.Contains(beforeChat.Text, sentence) {
						var chatId *uint = &beforeChat.Id
						result.ZoomChatHistory[index-1].ReactId = chatId
						break
					}

					searchCursor--
				}

				result.ZoomChatHistory[index-1].Text = emoji

			} else if category == NormalText {
				result.ZoomChatHistory[index-1].Text += text
			} else if category == Reply {
			}

			textCursor++

		}

	}

	return
}

func categorizeMessage(message string) ChatContentType {
	switch {
	case strings.HasPrefix(message, "Reacted to "):
		return Reaction

	case strings.HasPrefix(message, "Replying to "):
		return Reply

	case strings.HasPrefix(message, "Removed a "):
		return Remove

	default:
		return NormalText
	}
}

func extractReaction(message string) (emoji string, sentence string) {
	prefix := `Reacted to "`
	suffix := `" with `

	start := strings.Index(message, prefix) + len(prefix)
	end := strings.Index(message, suffix)

	sentence = message[start:end]
	sentence = strings.TrimRight(sentence, ".")

	emoji = message[end+len(suffix):]

	return
}

func extractReply(message string) (sentence string) {
	prefix := `Replying to "`

	start := strings.Index(message, prefix) + len(prefix)

	sentence = message[start : len(message)-1]
	return sentence
}