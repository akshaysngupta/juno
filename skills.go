package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var skillSet = map[string]*skillInput{
	"read-file": {
		Skill:       "read-file",
		Description: "This skill allows you to read a file.",
		Args: []string{
			"file-name",
		},
		Run: func(args []string) (*skillOutput, error) {
			filePath := args[0]
			output := readFile(repoRoot + "/" + filePath)
			return &skillOutput{
				Output: output,
			}, nil
		},
	},
	"write-file": {
		Skill:       "write-file",
		Description: "This skill writes to a file.",
		Args: []string{
			"file-name",
			"file-content",
		},
		Run: func(args []string) (*skillOutput, error) {
			filePath := args[0]
			content := args[1]
			writeFile(repoRoot+"/"+filePath, content)
			return &skillOutput{
				Output: "File written successfully.",
			}, nil
		},
	},
	"delete-file": {
		Skill:       "delete-file",
		Description: "This skill deletes a file.",
		Args: []string{
			"file-name",
		},
		Run: func(args []string) (*skillOutput, error) {
			filePath := args[0]
			deleteFile(repoRoot + "/" + filePath)
			return &skillOutput{
				Output: "File deleted successfully.",
			}, nil
		},
	},
	"list-files": {
		Skill:       "list-files",
		Description: "This skill lists the files.",
		Args:        []string{},
		Run: func(args []string) (*skillOutput, error) {
			output := getFolderStructure(repoRoot)
			return &skillOutput{
				Output: output,
			}, nil
		},
	},
	"done": {
		Skill:       "done",
		Description: "This skill indicates that work is finished.",
		Args:        []string{},
		Run: func(args []string) (*skillOutput, error) {
			return nil, nil
		},
	},
}

func parseSkill(input string) (*skillInput, error) {
	endPos := strings.Index(input, skillEndTag)

	if endPos == -1 {
		return nil, fmt.Errorf("skill tags not found. Input: %s", input)
	}

	skillStr := input[0:endPos]

	fmt.Println("Extraced skill", skillStr)

	var skill skillInput
	err := json.Unmarshal([]byte(skillStr), &skill)
	if err != nil {
		return nil, err
	}

	return &skill, nil
}

func performSkill(skill *skillInput) (*skillOutput, error) {
	fmt.Println("Working on skill", skill)
	skillDefinition, exists := skillSet[skill.Skill]
	if !exists {
		return nil, fmt.Errorf("command %v not recognized", skillDefinition)
	}

	skill.Run = skillDefinition.Run
	return skill.RunSkill()
}

func formatSkillOutput(skill *skillInput, output *skillOutput) string {
	skillBytes, _ := json.MarshalIndent(skill, "", "    ")
	resultBytes, _ := json.MarshalIndent(output, "", "    ")
	return fmt.Sprintf(formatterResult, skillBytes, resultBytes)
}

func initSkill() *skillInput {
	skill := *skillSet["list-files"]
	skill.Description = ""
	return &skill
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
