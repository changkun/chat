// Copyright 2023 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package openai_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"unicode/utf8"

	"changkun.de/x/chat/internal/openai"
)

func countWords(s string) int {
	counter := 0
	for _, word := range strings.Fields(s) {
		runeCount := utf8.RuneCountInString(word)
		if len(word) == runeCount {
			counter++
		} else {
			counter += runeCount
		}
	}

	return counter
}

func TestEdit(t *testing.T) {
	txt, err := os.ReadFile("testdata/data.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("total words: ", countWords(string(txt)))

	// TLDR summarization: https://beta.openai.com/examples/default-tldr-summary
	o, err := openai.Edit(context.Background(), &openai.EditRequest{
		Model:       "text-davinci-edit-001",
		Input:       string(txt),
		Instruction: "Fix the grammar and spelling mistakes.",
		Temperature: 0.5,
		N:           1,
	})
	if err != nil {
		t.Fatalf("failed to edit: %v", err)
	}

	for _, c := range o.Choices {
		t.Log(c.Text)
		t.Log("words: ", countWords(c.Text))
	}
}

func TestChat(t *testing.T) {
	respCh, errCh := openai.Chat(context.Background(), &openai.ChatRequest{
		Model:  "gpt-3.5-turbo-0301",
		Stream: false,
		Message: []openai.ChatMessage{
			{
				Role:    "system",
				Content: "You are a helpful assistant.",
			},
			{
				Role:    "user",
				Content: "What can you do?",
			},
		},
	})

	t.Log(<-respCh)
	t.Log(<-errCh)
}
