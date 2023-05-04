package chat

import (
	"bufio"
	"errors"
	"io"
	"math"
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

func (p *Parser) Parse(file io.Reader) (result Result, err error) {
	texts := make([]string, 0)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		texts = append(texts, scanner.Text())
	}

	metaIds := make(map[int]bool, 0)
	// find meta index
	for index, v := range texts {
		if !strings.Contains(v, "\t") {
			metaIds[index] = true
		}
	}

	// append extracted meta data with raw content
	var id int
	contents := make([]string, len(metaIds))
	var textCount int
	for i, v := range texts {
		if metaIds[i] {
			splited := strings.Split(v, " ")

			result.ZoomChatHistory = append(result.ZoomChatHistory, &ZoomChatHistory{
				Id:           uint(id),
				ChatedAt:     splited[0],
				SenderName:   splited[2],
				ReceiverName: splited[4],
			})

			id += 1
			textCount = 0
			continue
		}
		if textCount > 0 {
			contents[id-1] += "\n"
		}
		contents[id-1] += v
		textCount++
	}

	// append extracted raw content with Category
	for id, content := range contents {

		content = strings.TrimLeft(content, "\t")
		category := p.categorizeMessage(content)
		result.ZoomChatHistory[id].TextType = category

		if category == Reaction {
			emoji, sentence := p.extractReaction(content)
			parentId := p.findIDFromChatHistoryByText(sentence, uint(id), result.ZoomChatHistory)
			result.ZoomChatHistory[id].Text = emoji

			if parentId == nil {
				result.Statistic.MissingReactionIds = append(result.Statistic.MissingReactionIds, uint(id))
				continue
			}

			result.ZoomChatHistory[*parentId].ReactIds = append(result.ZoomChatHistory[*parentId].ReactIds, uint(id))

		} else if category == Reply {
			sentence, replyText := p.extractReply(content)
			parentId := p.findIDFromChatHistoryByText(sentence, uint(id), result.ZoomChatHistory)
			result.ZoomChatHistory[id].Text = replyText

			if parentId == nil {
				result.Statistic.MissingReplyIds = append(result.Statistic.MissingReplyIds, uint(id))
				continue
			}
			result.ZoomChatHistory[*parentId].ReplyIds = append(result.ZoomChatHistory[*parentId].ReplyIds, uint(id))

		} else if category == Remove {

			emoji, sentence := p.extractRemove(content)

			parentId := p.findIDFromChatHistoryByText(sentence, uint(id), result.ZoomChatHistory)

			result.ZoomChatHistory[id].Text = emoji
			result.ZoomChatHistory[id].Removed = true

			if parentId == nil {
				result.Statistic.MissingRemoveIds = append(result.Statistic.MissingRemoveIds, uint(id))
				continue
			}

			for index, v := range result.ZoomChatHistory[id].ReactIds {
				if v == uint(id) {
					result.ZoomChatHistory[id].ReactIds = append(result.ZoomChatHistory[id].ReactIds[:index], result.ZoomChatHistory[id].ReactIds[index+1:]...)
					break
				}
			}

		} else if category == NormalText {
			result.ZoomChatHistory[id].Text = content
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
	sentence = strings.Trim(sentence, `"`)
	emoji = message[end+len(suffix):]

	return
}

func (p *Parser) extractReply(message string) (sentence, replyText string) {
	prefix := `Replying to "`
	start := strings.Index(message, prefix) + len(prefix)

	end := strings.LastIndex(message, `"`)
	sentence = message[start:end]
	replyTextRaw := message[end+1:]

	sentence = strings.TrimRight(sentence, ".")

	for i, v := range strings.Split(replyTextRaw, "\n") {
		if i < 1 {
			continue
		}
		replyText += strings.TrimLeft(v, "\t")

		if i > 2 {
			replyText += "\n"
		}
	}

	return sentence, replyText
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
	cursor := currentIndex - 1
	count := 30
	for cursor < math.MaxUint32 && count >= 0 {

		beforeChat := history[cursor]
		if strings.Contains(text, "\t") {
			text = strings.Split(text, "\t")[0]
			text = strings.TrimRight(text, " ")
		}
		if strings.Contains(beforeChat.Text, text) {

			id = &beforeChat.Id
			return
		}

		cursor--
		count--
	}

	return
}
