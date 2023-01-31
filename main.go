package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/containrrr/shoutrrr"
	"github.com/google/go-github/v50/github"
)

func main() {
	issuesToSend, err := getIssuesToSend()

	if err != nil {
		log.Fatal(err)
	}

	if len(issuesToSend) == 0 {
		fmt.Println("No issues to send")
		return
	}

	err = sendIssuesViaMail(issuesToSend)

	if err != nil {
		log.Fatal(err)
	}

	for _, issue := range issuesToSend {
		err = appendIssue(issue)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func sendIssuesViaMail(issues []*github.Issue) error {
	formatted := ""

	for _, issue := range issues {
		formatted += fmt.Sprintf("%d. %v - %v\n", *issue.Number, *issue.Title, issue.GetHTMLURL())
	}

	url := os.Getenv("SHOUTRRR_URL")

	return shoutrrr.Send(url, formatted)
}

func getIssuesToSend() ([]*github.Issue, error) {
	client := github.NewClient(nil)

	issues, _, err := client.Issues.ListByRepo(context.Background(), "elk-zone", "elk", &github.IssueListByRepoOptions{Labels: []string{"good first issue"}})

	if err != nil {
		return nil, err
	}

	issuesToSend, err := filterSentIssues(issues)

	if err != nil {
		return nil, err
	}

	return issuesToSend, nil
}

func filterSentIssues(issues []*github.Issue) ([]*github.Issue, error) {
	content, err := os.ReadFile("issues.txt")

	if err != nil && errors.Is(err, fs.ErrNotExist) {
		fmt.Print("Creating issues.txt file... ")
		err = os.WriteFile("issues.txt", []byte(""), 0644)
	}

	if err != nil {
		return nil, err
	}

	readIds := strings.Split(string(content), ",")

	issuesToSend := []*github.Issue{}

	for _, issue := range issues {
		found := false
		for _, line := range readIds {
			if strconv.Itoa(*issue.Number) == line {
				found = true
			}
		}
		if !found {
			issuesToSend = append(issuesToSend, issue)
		}
	}
	return issuesToSend, nil
}

func appendIssue(issue *github.Issue) error {
	content, err := os.ReadFile("issues.txt")

	if err != nil {
		return err
	}

	content = append(content, []byte(strconv.Itoa(*issue.Number)+",")...)

	err = os.WriteFile("issues.txt", content, 0644)

	if err != nil {
		return err
	}

	return nil
}
