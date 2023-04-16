package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

const (
	repoRoot = "/Users/akshaysngupta/Documents/juno/helloworld"
)

func executeTask(client *openai.Client, systemCard string, initTaskCard string) {
	action := initAction()
	output, _ := performAction(action)
	initTaskCard = initTaskCard + formatActionOutput(action, output)

	reset := true
	var updatedTaskCard string
	var request []openai.ChatCompletionMessage
	for {
		if reset {
			updatedTaskCard = initTaskCard
			writeFile("request.txt", updatedTaskCard)
			reset = false
		}

		request = createRequest(updatedTaskCard)
		fmt.Println("Sending Juno Request:\n", messagesToFormattedjson(request))
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:       openai.GPT3Dot5Turbo,
				Messages:    request,
				Temperature: 0.9,
			},
		)
		if err != nil {
			panic(err)
		}

		response := resp.Choices[0].Message.Content
		fmt.Println("Juno Response:\n", response)

		action, err = parseAction(response)
		if err != nil {
			fmt.Println("Resetting due to error", err)
			reset = true
			time.Sleep(10 * time.Second)
			continue
		}

		output, err := performAction(action)
		if err != nil {
			fmt.Println("Resetting due to error", err)
			reset = true
			time.Sleep(10 * time.Second)
			continue
		}

		if output == nil {
			fmt.Println("Exiting.")
			return
		}

		updatedTaskCard = updatedTaskCard + formatActionOutput(action, output)
		writeFile("request.txt", updatedTaskCard)
	}
}

func messagesToFormattedjson(messages []openai.ChatCompletionMessage) string {
	marshaled, _ := json.MarshalIndent(messages, "", "  ")
	return string(marshaled)
}

func createRequest(task string) []openai.ChatCompletionMessage {
	return []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: string(systemCard),
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: string(task),
		},
	}
}
