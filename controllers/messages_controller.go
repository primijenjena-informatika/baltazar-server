package controllers

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/openai/openai-go"
	"github.com/primijenjan-informatika/baltazar-server/llm"
	"github.com/spf13/viper"
)

type MessageRole string

type MessageContent string

type Model string

type Message struct {
	Role    MessageRole    `json:"role"`
	Content MessageContent `json:"content"`
}

type SendMessageRequest struct {
	Messages []Message `json:"messages"`
	Model    Model
}

func HandleMessages(c echo.Context) error {
	apiKeys := viper.GetStringSlice("api.api_keys")

	auth := c.Request().Header.Get("Authorization")
	apiKey := strings.TrimPrefix(auth, "Bearer ")

	if !slices.Contains(apiKeys, apiKey) {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "Invalid API key",
		})
	}

	req := new(SendMessageRequest)

	err := c.Bind(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Invalid request schema",
		})
	}

	if req.Model != Model(viper.GetString("model.id")) {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Invalid model id",
		})
	}

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	tokenChan := make(chan string)

	go func() {
		defer close(tokenChan)

		stream := llm.SendChatContext(createChatCompletionMessageParams(req.Messages))
		for stream.Next() {
			token := stream.Current().Choices[0].Delta.Content
			tokenChan <- token
		}
		if stream.Err() != nil {
			c.Logger().Error(stream.Err())
		}
	}()

	c.Response().Writer.WriteHeader(http.StatusOK)
	for token := range tokenChan {
		_, err := fmt.Fprintf(c.Response().Writer, "data: %s\n\n", token)
		if err != nil {
			c.Logger().Error(err)
			break
		}
		c.Response().Flush()
	}

	return nil
}

func createChatCompletionMessageParams(messages []Message) []openai.ChatCompletionMessageParamUnion {
	var chatCompletionMessages []openai.ChatCompletionMessageParamUnion

	for _, message := range messages {
		if message.Role != "user" && message.Role != "assistant" {
			continue
		}

		chatCompletionMessages = append(chatCompletionMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatCompletionMessageRole(string(message.Role)),
			Content: string(message.Content),
		})
	}

	return chatCompletionMessages
}
