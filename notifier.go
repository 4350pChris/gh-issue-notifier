package main

import (
	"fmt"
	"log"

	"github.com/containrrr/shoutrrr"
)

func NotifyOfIssues(url string, repos []*Repository) error {
	formatted := ""

	for _, repo := range repos {
		formatted += fmt.Sprintf("Repo: %v\n", repo.FullName)
		for _, issue := range repo.Issues {
			formatted += fmt.Sprintf("%d. %v - %v\n", issue.Number, issue.Title, issue.HtmlUrl)
		}
		formatted += "\n"
	}

	log.Default().Println("Sending notification...")

	return shoutrrr.Send(url, formatted)
}
