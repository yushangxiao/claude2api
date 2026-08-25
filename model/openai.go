package model

import (
	"claude2api/logger"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatCompletionRequest struct {
	Model    string                   `json:"model"`
	Messages []map[string]interface{} `json:"messages"`
	Stream   bool                     `json:"stream"`
	Tools    []map[string]interface{} `json:"tools,omitempty"`
}

// OpenAISrteamResponse 定义 OpenAI 的流式响应结构
type OpenAISrteamResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []StreamChoice `json:"choices"`
}

// Choice 结构表示 OpenAI 返回的单个选项
type StreamChoice struct {
	Index        int         `json:"index"`
	Delta        Delta       `json:"delta"`
	Logprobs     interface{} `json:"logprobs"`
	FinishReason interface{} `json:"finish_reason"`
}

type NoStreamChoice struct {
	Index        int         `json:"index"`
	Message      Message     `json:"message"`
	Logprobs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
}

// Delta 结构用于存储返回的文本内容
type Delta struct {
	Content string `json:"content"`
}
type Message struct {
	Role       string        `json:"role"`
	Content    string        `json:"content"`
	Refusal    interface{}   `json:"refusal"`
	Annotation []interface{} `json:"annotation"`
}

type OpenAIResponse struct {
	ID      string           `json:"id"`
	Object  string           `json:"object"`
	Created int64            `json:"created"`
	Model   string           `json:"model"`
	Choices []NoStreamChoice `json:"choices"`
	Usage   Usage            `json:"usage"`
}
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func ReturnOpenAIResponse(text string, stream bool, gc *gin.Context, modelName string) error {
	if modelName == "" {
		modelName = "claude-3-7-sonnet-20250219"
	}
	if stream {
		return streamRespose(text, modelName, gc)
	} else {
		return noStreamResponse(text, modelName, gc)
	}
}

func streamRespose(text string, modelName string, gc *gin.Context) error {
	openAIResp := &OpenAISrteamResponse{
		ID:      uuid.New().String(),
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   modelName,
		Choices: []StreamChoice{
			{
				Index: 0,
				Delta: Delta{
					Content: text,
				},
				Logprobs:     nil,
				FinishReason: nil,
			},
		},
	}

	jsonBytes, err := json.Marshal(openAIResp)
	jsonBytes = append([]byte("data: "), jsonBytes...)
	jsonBytes = append(jsonBytes, []byte("\n\n")...)
	if err != nil {
		logger.Error(fmt.Sprintf("Error marshalling JSON: %v", err))
		return err
	}

	// 发送数据
	gc.Writer.Write(jsonBytes)
	gc.Writer.Flush()
	return nil
}

func noStreamResponse(text string, modelName string, gc *gin.Context) error {
	openAIResp := &OpenAIResponse{
		ID:      uuid.New().String(),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   modelName,
		Choices: []NoStreamChoice{
			{
				Index: 0,
				Message: Message{
					Role:    "assistant",
					Content: text,
				},
				Logprobs:     nil,
				FinishReason: "stop",
			},
		},
	}

	gc.JSON(200, openAIResp)
	return nil
}
