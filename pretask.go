package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func populateActions(preTaskCard string) string {
	var actionCards string
	for _, action := range actionSet {
		actionBytes, _ := json.MarshalIndent(action, "", "    ")
		actionCard := fmt.Sprintf(
			"Action: %s\n"+
				"Description: %s\n"+
				"Example:\n%s\n\n",
			action.Action,
			action.Description,
			string(actionBytes))
		actionCards += actionCard
	}

	return strings.Replace(preTaskCard, "${ACTIONS}", actionCards, 1)
}
