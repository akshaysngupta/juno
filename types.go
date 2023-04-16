package main

import "encoding/json"

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

	Description string `json:"description,omitempty"`

	Args []string `json:"args,omitempty"`

	Run func(args []string) (*actionOutput, error)
}

func (action *actionInput) RunAction() (*actionOutput, error) {
	return action.Run(action.Args)
}

// Marshal returns the JSON encoding of the actionInput
func (a *actionInput) MarshalJSON() ([]byte, error) {
	// create a map to store the fields of the actionInput
	m := make(map[string]interface{})
	m["action"] = a.Action
	if a.Description != "" {
		m["description"] = a.Description
	}
	if len(a.Args) > 0 {
		m["args"] = a.Args
	}
	// use json.Marshal to encode the map
	return json.Marshal(m)
}

type actionOutput struct {
	Output string `json:"output"`
}
