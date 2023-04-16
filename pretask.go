package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func populateSkills(preTaskCard string) string {
	var skillCards string
	for _, skill := range skillSet {
		skillBytes, _ := json.MarshalIndent(skill, "", "    ")
		skillCard := fmt.Sprintf(
			"Skill: %s\n"+
				"Description: %s\n"+
				"Example:\n%s\n\n",
			skill.Skill,
			skill.Description,
			string(skillBytes))
		skillCards += skillCard
	}

	return strings.Replace(preTaskCard, "${SKILLS}", skillCards, 1)
}
