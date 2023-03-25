package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var (
	actionStartTag = "<action-start>"

	actionEndTag = "<action-end>"

	resultStartTag = "<result-start>"

	resultEndTag = "<result-end>"

	formatterResult = "%s\n" + actionEndTag + "\n" +
		resultStartTag + "\n%s\n" + resultEndTag + "\n" +
		actionStartTag + "\n"
)

type actionInput struct {
	Action string `json:"action"`

	Args []string `json:"args,omitempty"`
}

type actionOutput struct {
	Output string `json:"output"`
}

func parseAction(input string) (*actionInput, error) {
	endPos := strings.Index(input, actionEndTag)

	if endPos == -1 {
		return nil, fmt.Errorf("action tags not found. Input: %s", input)
	}

	actionStr := input[0:endPos]

	fmt.Println("Extraced Action", actionStr)

	var action actionInput
	err := json.Unmarshal([]byte(actionStr), &action)
	if err != nil {
		return nil, err
	}

	return &action, nil
}

func performAction(action *actionInput) (*actionOutput, error) {
	fmt.Println("Working on Action", action)
	switch action.Action {
	case "read-file":
		filePath := action.Args[0]
		output := readFile(repoRoot + "/" + filePath)
		return &actionOutput{
			Output: output,
		}, nil
	case "delete-file":
		filePath := action.Args[0]
		deleteFile(repoRoot + "/" + filePath)
		return &actionOutput{
			Output: "File deleted successfully.",
		}, nil
	case "write-file":
		filePath := action.Args[0]
		content := action.Args[1]
		writeFile(repoRoot+"/"+filePath, content)
		return &actionOutput{
			Output: "File written successfully.",
		}, nil
	case "list-files":
		output := getFolderStructure(repoRoot)
		return &actionOutput{
			Output: output,
		}, nil
	case "quit":
		return nil, nil
	}

	return nil, fmt.Errorf("command %s not recognized", action.Action)
}

func formatActionOutput(action *actionInput, output *actionOutput) string {
	actionBytes, _ := json.MarshalIndent(action, "", "    ")
	resultBytes, _ := json.MarshalIndent(output, "", "    ")
	return fmt.Sprintf(formatterResult, actionBytes, resultBytes)
}

func initAction() *actionInput {
	return &actionInput{
		Action: "list-files",
	}
}

func getFolderStructure(path string) string {
	files, err := os.ReadDir(path)
	if err != nil {
		return err.Error()
	}

	fileNames := make([]string, 0)
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	return strings.Join(fileNames, ", ")
}

func readFile(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		return err.Error()
	}

	return string(content)
}

func writeFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), os.ModePerm)
}

func deleteFile(path string) {
	_ = os.Remove(path)
}
