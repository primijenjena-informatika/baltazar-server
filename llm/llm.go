package llm

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/spf13/viper"
)

func SendChatContext(messages []openai.ChatCompletionMessageParamUnion) {
	client := openai.NewClient(
		option.WithAPIKey(viper.GetString("model.api_key")),
		option.WithBaseURL(viper.GetString("model.api_base_url")),
	)

	ctx := context.Background()

	stream := client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F(append([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(viper.GetString("model.system_prompt")),
		},
			messages...,
		)),
		Model: openai.F(openai.ChatModel(viper.GetString("model.id"))),
	})

}
