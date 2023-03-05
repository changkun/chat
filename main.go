// Copyright 2023 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"changkun.de/x/chat/internal/openai"
	"changkun.de/x/chat/internal/term"
)

func main() {
	if os.Getenv("OPENAI_API_KEY") == "" {
		fmt.Fprint(os.Stderr, "Please set OPENAI_API_KEY environment variable.\n")
		return
	}

	stdin := os.Stdin
	stdout := os.Stdout
	fmt.Fprint(stdout, term.Orange("Hi, I'm a chatbot. How can I help you?\n"))

	session := []openai.ChatMessage{
		{
			Role:    "system",
			Content: "You are a helpful assistant.",
		},
	}

	for {
		fmt.Fprint(stdout, term.Orange("User: "))
		buf := bytes.NewBuffer(nil)
		_, err := io.Copy(buf, stdin)
		if err != nil {
			fmt.Fprintf(stdout, "Error: %v", err)
			return
		}

		userMsg := openai.ChatMessage{
			Role:    "user",
			Content: buf.String(),
		}
		session = append(session, userMsg)

		respCh, errCh := openai.Chat(context.Background(), &openai.ChatRequest{
			Model:   "gpt-3.5-turbo-0301",
			Stream:  true,
			Message: session,
		})

		response := openai.ChatMessage{
			Role:    "assistant",
			Content: "",
		}
		fmt.Fprint(stdout, term.Orange("Assistant: "))

	streamLoop:
		for {
			select {
			case r, ok := <-respCh:
				if !ok {
					break streamLoop
				}
				if len(r.Choices) == 0 {
					fmt.Fprint(stdout, "No response. End of the session.")
					return
				}

				msg := r.Choices[0].Delta
				fmt.Fprint(stdout, msg.Content)
				response.Content += msg.Content
			case err, ok := <-errCh:
				if !ok {
					if err != nil {
						fmt.Fprintf(stdout, "Error: %v", err)
					}
					break streamLoop
				}
			}
		}
		session = append(session, response)
		fmt.Fprintf(stdout, "\n")
	}
}
