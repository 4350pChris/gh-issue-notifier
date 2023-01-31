package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
	"time"

	"github.com/containrrr/shoutrrr"
	"github.com/google/go-github/v50/github"
	_ "github.com/joho/godotenv/autoload"
)

type Sendable struct {
	repo   string
	issues []*github.Issue
}

func main() {
	for {
		err := doWork()
		if err != nil {
			log.Fatal(err)
		}

		duration, err := time.ParseDuration(os.Getenv("INTERVAL"))

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Sleeping for %v\n", duration)
		time.Sleep(duration)
	}
}

func doWork() error {
	sendables, err := getIssuesToSend()

	if err != nil {
		return err
	}

	if len(sendables) == 0 {
		fmt.Println("No issues to send")
		return nil
	}

	err = sendIssuesViaMail(sendables)

	if err != nil {
		return err
	}

	// save all issues to txt file
	for _, sendable := range sendables {
		for _, issue := range sendable.issues {
			err = appendIssue(issue)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func sendIssuesViaMail(sendables []Sendable) error {
	formatted := ""

	for _, sendable := range sendables {
		formatted += fmt.Sprintf("Repo: %v\n", sendable.repo)
		for _, issue := range sendable.issues {
			formatted += fmt.Sprintf("%d. %v - %v\n", *issue.Number, *issue.Title, issue.GetHTMLURL())
		}
		formatted += "\n"
	}

	url := os.Getenv("SHOUTRRR_URL")

	return shoutrrr.Send(url, formatted)
}

func getIssuesToSend() ([]Sendable, error) {
	client := github.NewClient(nil)

	patterns := strings.Split(os.Getenv("WATCH_PATTERNS"), ",")

	sendables := []Sendable{}

	for _, pattern := range patterns {

		splitPattern := strings.Split(pattern, "/")

		if len(splitPattern) != 3 {
			return nil, fmt.Errorf("invalid watch pattern %v", pattern)
		}

		owner, repo, label := splitPattern[0], splitPattern[1], splitPattern[2]

		issues, _, err := client.Issues.ListByRepo(context.Background(), owner, repo, &github.IssueListByRepoOptions{Labels: []string{label}})
		if err != nil {
			return nil, err
		}

		issuesToSend, err := filterSentIssues(issues)
		if err != nil {
			return nil, err
		}

		// skip repos that have no new issues
		if len(issuesToSend) > 0 {
			sendables = append(sendables, Sendable{repo: repo, issues: issuesToSend})
		}
	}

	return sendables, nil
}

func filterSentIssues(issues []*github.Issue) ([]*github.Issue, error) {
	content, err := os.ReadFile("issues.txt")

	if err != nil && errors.Is(err, fs.ErrNotExist) {
		fmt.Println("Creating issues.txt file... ")
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
			if *issue.URL == line {
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

	content = append(content, []byte(*issue.URL+",")...)

	err = os.WriteFile("issues.txt", content, 0644)

	if err != nil {
		return err
	}

	return nil
}
