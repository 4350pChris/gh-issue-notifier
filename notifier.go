package main

import (
	"fmt"

	"github.com/containrrr/shoutrrr"
)

func NotifyOfIssues(url string, sendables []Sendable) error {
	formatted := ""

	for _, sendable := range sendables {
		formatted += fmt.Sprintf("Repo: %v\n", sendable.repo)
		for _, issue := range sendable.issues {
			formatted += fmt.Sprintf("%d. %v - %v\n", *issue.Number, *issue.Title, issue.GetHTMLURL())
		}
		formatted += "\n"
	}

	return shoutrrr.Send(url, formatted)
}
