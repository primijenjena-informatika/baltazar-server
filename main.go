package main

import (
	"fmt"

	"github.com/openai/openai-go"
	"github.com/primijenjan-informatika/baltazar-server/config"
	"github.com/primijenjan-informatika/baltazar-server/llm"
)

func main() {
	config.InitConfig()
	stream := llm.SendChatContext([]openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("Kako da napravim web stranicu?"),
	})

	for stream.Next() {
		fmt.Print(stream.Current().Choices[0].Delta.Content)
	}
	fmt.Print("\n")
}
