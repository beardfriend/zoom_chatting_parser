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

	ReactIds []uint

	ReplyIds []uint

	TextType ChatContentType

	SenderName string

	ReceiverName string

	Text string

	Removed bool

	// format HH:MM:ss
	ChatedAt string
}

type Result struct {
	ZoomChatHistory []*ZoomChatHistory
	Statistic       struct {
		MissingReactionIds []uint
		MissingReplyIds    []uint
		MissingRemoveIds   []uint
	}
}

var ErrorNoFile = errors.New("no such file")

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(file *os.File) (result Result, err error) {
	if file == nil {
		err = ErrorNoFile
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var Id, textCursor uint
	var reactionText string
	for scanner.Scan() {
		line := scanner.Text()

		isMetaData := !strings.Contains(line, "\t")
		// MetaData
		if isMetaData {
			splited := strings.Split(line, " ")

			result.ZoomChatHistory = append(result.ZoomChatHistory, &ZoomChatHistory{
				Id:           Id,
				ChatedAt:     splited[0],
				SenderName:   splited[2],
				ReceiverName: splited[4],
			})

			Id++
			textCursor = 0
			continue
		}

		// Text Zone

		text := strings.TrimPrefix(line, "\t")
		category := p.categorizeMessage(text)
		result.ZoomChatHistory[Id-1].TextType = category
		// 리엑션일 때는 리엑션이 어느 글에서 했는지 판단
		if category == Reaction {
			// 리엑션이 띄어쓰기가 됐을 경우
			reactionText = text
			if !strings.Contains(text, `" with `) {
				reactionText += text
				continue
			}

			emoji, sentence := p.extractReaction(reactionText)
			id := p.findIDFromChatHistoryByText(sentence, Id, result.ZoomChatHistory)
			result.ZoomChatHistory[Id-1].Text = emoji

			if id == nil {
				result.Statistic.MissingReactionIds = append(result.Statistic.MissingReactionIds, Id)
				continue
			}

			result.ZoomChatHistory[*id].ReactIds = append(result.ZoomChatHistory[*id].ReactIds, Id)

		} else if category == Reply {

			sentence := p.extractReply(text)
			id := p.findIDFromChatHistoryByText(sentence, Id, result.ZoomChatHistory)
			result.ZoomChatHistory[Id-1].Text = sentence
			if id == nil {
				result.Statistic.MissingReplyIds = append(result.Statistic.MissingReplyIds, Id)
				continue
			}
			result.ZoomChatHistory[*id].ReplyIds = append(result.ZoomChatHistory[*id].ReplyIds, Id)

		} else if category == Remove {

			emoji, sentence := p.extractRemove(text)

			id := p.findIDFromChatHistoryByText(sentence, Id, result.ZoomChatHistory)
			result.ZoomChatHistory[Id-1].Text = emoji
			result.ZoomChatHistory[Id-1].Removed = true
			if id == nil {
				result.Statistic.MissingRemoveIds = append(result.Statistic.MissingRemoveIds, Id)
				continue
			}

			for index, v := range result.ZoomChatHistory[*id].ReactIds {
				if v == Id {
					result.ZoomChatHistory[*id].ReactIds = append(result.ZoomChatHistory[*id].ReactIds[:index], result.ZoomChatHistory[*id].ReactIds[index+1:]...)
					break
				}
			}

		} else if category == NormalText {
			if textCursor > 0 {
				result.ZoomChatHistory[Id-1].Text += "\n"
			}
			result.ZoomChatHistory[Id-1].Text += text

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
