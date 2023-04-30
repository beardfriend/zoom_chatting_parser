package chat

import (
	"bufio"
	"errors"
	"math"
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

	var index, textCursor uint

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
			category := p.categorizeMessage(text)

			// 리엑션일 때는 리엑션이 어느 글에서 했는지 판단
			if category == Reaction {
				emoji, sentence := p.extractReaction(text)

				id := p.findIDFromChatHistoryByText(sentence, index, result.ZoomChatHistory)
				result.ZoomChatHistory[index-1].ReactId = id
				result.ZoomChatHistory[index-1].Text = emoji

			} else if category == NormalText {
				result.ZoomChatHistory[index-1].Text += text
			} else if category == Reply {
				// sentence := p.extractReply(text)
			}

			textCursor++

		}

	}

	return
}

func (p *Parser) categorizeMessage(message string) ChatContentType {
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

func (p *Parser) extractReaction(message string) (emoji string, sentence string) {
	prefix := `Reacted to "`
	suffix := `" with `

	start := strings.Index(message, prefix) + len(prefix)
	end := strings.Index(message, suffix)

	sentence = message[start:end]
	sentence = strings.TrimRight(sentence, ".")

	emoji = message[end+len(suffix):]

	return
}

func (p *Parser) extractReply(message string) (sentence string) {
	prefix := `Replying to "`

	start := strings.Index(message, prefix) + len(prefix)

	sentence = message[start : len(message)-1]
	sentence = strings.TrimRight(sentence, ".")
	return sentence
}

func (p *Parser) extractRemove(message string) (emoji, sentence string) {
	prefix := `Removed a `
	suffix := " reaction from"
	start := strings.Index(message, prefix) + len(prefix)
	end := strings.Index(message, suffix)

	emoji = message[start:end]

	sentence = message[end+len(suffix)+2 : len(message)-2]
	sentence = strings.TrimRight(sentence, ".")
	return
}

func (p *Parser) findIDFromChatHistoryByText(text string, currentIndex uint, history []*ZoomChatHistory) (id *uint) {
	cursor := currentIndex - 2
	count := 10
	for cursor < math.MaxUint32 && count >= 0 {

		beforeChat := history[cursor]

		if strings.Contains(beforeChat.Text, text) {
			id = &beforeChat.Id
			return
		}

		cursor--
		count--
	}

	return
}
