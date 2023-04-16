package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var actionSet = map[string]*actionInput{
	"read-file": {
		Action:      "read-file",
		Description: "This action read a file.",
		Args: []string{
			"file-name",
		},
		Run: func(args []string) (*actionOutput, error) {
			filePath := args[0]
			output := readFile(repoRoot + "/" + filePath)
			return &actionOutput{
				Output: output,
			}, nil
		},
	},
	"write-file": {
		Action:      "write-file",
		Description: "This action writes to a file.",
		Args: []string{
			"file-name",
			"file-content",
		},
		Run: func(args []string) (*actionOutput, error) {
			filePath := args[0]
			content := args[1]
			writeFile(repoRoot+"/"+filePath, content)
			return &actionOutput{
				Output: "File written successfully.",
			}, nil
		},
	},
	"delete-file": {
		Action:      "delete-file",
		Description: "This action deletes a file.",
		Args: []string{
			"file-name",
		},
		Run: func(args []string) (*actionOutput, error) {
			filePath := args[0]
			deleteFile(repoRoot + "/" + filePath)
			return &actionOutput{
				Output: "File deleted successfully.",
			}, nil
		},
	},
	"list-files": {
		Action:      "list-files",
		Description: "This action lists the files.",
		Args:        []string{},
		Run: func(args []string) (*actionOutput, error) {
			output := getFolderStructure(repoRoot)
			return &actionOutput{
				Output: output,
			}, nil
		},
	},
	"done": {
		Action:      "done",
		Description: "This action indicates that work is finished.",
		Args:        []string{},
		Run: func(args []string) (*actionOutput, error) {
			return nil, nil
		},
	},
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
	actionDefinition, exists := actionSet[action.Action]
	if !exists {
		return nil, fmt.Errorf("command %v not recognized", actionDefinition)
	}

	action.Run = actionDefinition.Run
	return action.RunAction()
}

func formatActionOutput(action *actionInput, output *actionOutput) string {
	actionBytes, _ := json.MarshalIndent(action, "", "    ")
	resultBytes, _ := json.MarshalIndent(output, "", "    ")
	return fmt.Sprintf(formatterResult, actionBytes, resultBytes)
}

func initAction() *actionInput {
	action := *actionSet["list-files"]
	action.Description = ""
	return &action
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
