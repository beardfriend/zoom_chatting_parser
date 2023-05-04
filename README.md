# zoom_chatting_parser

A library for golang to parse chat history in Zoom

## Quick Installation

```bash
go get "github.com/beardfriend/zoom_chatting_parser"
```

## Usage

1. Save chat history in Zoom.

![](docs/chat_save.png)

2. Open the Zoom chat file and put it into the parser.

```go
package main

import (
	"fmt"
	"os"

	z "github.com/beardfriend/zoom_chatting_parser"
)


func main() {
	file, _ := os.Open("meeting_saved_chat.txt")
	parser := z.NewParser()
	result, err := parser.Parse(file)
	if err != nil {
		panic(err)
	}
	for _, v := range result.ZoomChatHistory {
		fmt.Println(v.Id, v.ChatedAt, v.ReactIds, v.ReplyIds, v.ReceiverName, v.SenderName, v.Text, v.TextType, v.Removed)
	}
}

```

3. You can see all of zoom chat history.
