package main

import "encoding/json"

var (
	skillStartTag = "<skill-start>"

	skillEndTag = "<skill-end>"

	resultStartTag = "<result-start>"

	resultEndTag = "<result-end>"

	formatterResult = "%s\n" + skillEndTag + "\n" +
		resultStartTag + "\n%s\n" + resultEndTag + "\n" +
		skillStartTag + "\n"
)

type skillInput struct {
	Skill string `json:"skill"`

	Description string `json:"description,omitempty"`

	Args []string `json:"args,omitempty"`

	Run func(args []string) (*skillOutput, error)
}

func (skill *skillInput) RunSkill() (*skillOutput, error) {
	return skill.Run(skill.Args)
}

// Marshal returns the JSON encoding of the skillInput
func (a *skillInput) MarshalJSON() ([]byte, error) {
	// create a map to store the fields of the skillInput
	m := make(map[string]interface{})
	m["skill"] = a.Skill
	if a.Description != "" {
		m["description"] = a.Description
	}
	if len(a.Args) > 0 {
		m["args"] = a.Args
	}
	// use json.Marshal to encode the map
	return json.Marshal(m)
}

type skillOutput struct {
	Output string `json:"output"`
}
