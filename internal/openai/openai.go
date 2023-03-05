// Copyright 2023 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package openai

import (
	"context"
	"os"
	"time"
)

var (
	apiKey = os.Getenv("OPENAI_API_KEY")
)

type CompletionRequest struct {
	Model            string  `json:"model,omitempty"`
	Prompt           string  `json:"prompt,omitempty"`
	Suffix           string  `json:"suffix,omitempty"`
	MaxTokens        int     `json:"max_tokens,omitempty"`
	Temperature      float64 `json:"temperature,omitempty"`
	TopP             float64 `json:"top_p,omitempty"`
	N                int     `json:"n,omitempty"`
	Stream           bool    `json:"stream,omitempty"`
	LogProbs         int     `json:"logprobs,omitempty"`
	Echo             bool    `json:"echo,omitempty"`
	Stop             string  `json:"stop,omitempty"`
	PresencePenalty  float64 `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`
	BestOf           int     `json:"best_of,omitempty"`
}

type CompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string `json:"text"`
		Index        int    `json:"index"`
		LogProbs     int    `json:"logprobs"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int64 `json:"prompt_tokens"`
		CompletionTokens int64 `json:"completion_tokens"`
		TotalTokens      int64 `json:"total_tokens"`
	} `json:"usage"`
}

// Completion completes the input text based on the prompt.
//
// See: https://platform.openai.com/docs/api-reference/completions
func Completion(ctx context.Context, url string, in *CompletionRequest) (*CompletionResponse, error) {
	if !in.Stream {
		return httpRequest[CompletionRequest, CompletionResponse](ctx, "https://api.openai.com/v1/completions", in)
	}
	panic("not implemented")
}

type EditRequest struct {
	Model       string  `json:"model"`
	Input       string  `json:"input"`       // text
	Instruction string  `json:"instruction"` // prompt
	Temperature float64 `json:"temperature"` // creativity
	N           int     `json:"n"`
}

type EditResponse struct {
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Choices []struct {
		Text  string `json:"text"`
		Index int    `json:"index"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int64 `json:"prompt_tokens"`
		CompletionTokens int64 `json:"completion_tokens"`
		TotalTokens      int64 `json:"total_tokens"`
	} `json:"usage"`
}

type EditOutput struct {
	Choices []string
	Created time.Time
}

// Edit edits the input text based on the prompt.
//
// Recommended model is "text-davinci-edit-001".
// See: https://platform.openai.com/docs/api-reference/edits
func Edit(ctx context.Context, in *EditRequest) (*EditResponse, error) {
	return httpRequest[EditRequest, EditResponse](ctx, "https://api.openai.com/v1/edits", in)
}

type ChatRequest struct {
	Model            string        `json:"model,omitempty"`
	Message          []ChatMessage `json:"messages"`
	Temperature      float64       `json:"temperature,omitempty"`
	TopP             float64       `json:"top_p,omitempty"`
	N                int           `json:"n,omitempty"`
	Stream           bool          `json:"stream,omitempty"`
	Stop             []string      `json:"stop,omitempty"`
	MaxTokens        int           `json:"max_tokens,omitempty"`
	PresencePenalty  float64       `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64       `json:"frequency_penalty,omitempty"`
}

type ChatMessage struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Choices []struct {
		Index        int         `json:"index"`
		Delta        ChatMessage `json:"delta"`
		Message      ChatMessage `json:"message"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int64 `json:"prompt_tokens"`
		CompletionTokens int64 `json:"completion_tokens"`
		TotalTokens      int64 `json:"total_tokens"`
	} `json:"usage"`
}

// Chat creates a chat request to ChatGPT.
func Chat(ctx context.Context, in *ChatRequest) (<-chan *ChatResponse, <-chan error) {
	chatEndpoint := "https://api.openai.com/v1/chat/completions"
	ch := make(chan *ChatResponse)
	errCh := make(chan error)
	go func() {
		if !in.Stream {
			resp, err := httpRequest[ChatRequest, ChatResponse](ctx, chatEndpoint, in)
			ch <- resp
			errCh <- err
			close(ch)
			close(errCh)
		}

		errCh <- httpStream(ctx, chatEndpoint, in, ch)
		close(errCh)
	}()
	return ch, errCh
}
