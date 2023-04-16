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
	skill := initSkill()
	output, _ := performSkill(skill)
	initTaskCard = initTaskCard + formatSkillOutput(skill, output)

	reset := true
	var updatedTaskCard string
	var messages []openai.ChatCompletionMessage
	for {
		if reset {
			updatedTaskCard = initTaskCard
			writeFile("request.txt", updatedTaskCard)
			reset = false
		}

		messages = createMessages(updatedTaskCard)
		fmt.Println("Sending Juno Request:\n", messagesToFormattedjson(messages))
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:       openai.GPT3Dot5Turbo,
				Messages:    messages,
				Temperature: 0.9,
			},
		)
		if err != nil {
			panic(err)
		}

		response := resp.Choices[0].Message.Content
		fmt.Println("Juno Response:\n", response)

		skill, err = parseSkill(response)
		if err != nil {
			fmt.Println("Resetting due to error", err)
			reset = true
			time.Sleep(10 * time.Second)
			continue
		}

		output, err := performSkill(skill)
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

		updatedTaskCard = updatedTaskCard + formatSkillOutput(skill, output)
		writeFile("request.txt", updatedTaskCard)
	}
}

func messagesToFormattedjson(messages []openai.ChatCompletionMessage) string {
	marshaled, _ := json.MarshalIndent(messages, "", "  ")
	return string(marshaled)
}

func createMessages(task string) []openai.ChatCompletionMessage {
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
