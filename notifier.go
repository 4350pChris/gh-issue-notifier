package main

import (
	"fmt"
	"log"

	"github.com/containrrr/shoutrrr"
)

func NotifyOfIssues(url string, repos []*Repository) error {
	formatted := ""

	for _, repo := range repos {
		formatted += fmt.Sprintf("Repo: %v\n", repo.fullName)
		for _, issue := range repo.issues {
			formatted += fmt.Sprintf("%d. %v - %v\n", issue.number, issue.title, issue.htmlUrl)
		}
		formatted += "\n"
	}

	log.Default().Println("Sending notification...")

	return shoutrrr.Send(url, formatted)
}
