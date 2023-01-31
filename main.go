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
	issuesToSend, err := getIssuesToSend()

	if err != nil {
		return err
	}

	if len(issuesToSend) == 0 {
		fmt.Println("No issues to send")
		return nil
	}

	err = sendIssuesViaMail(issuesToSend)

	if err != nil {
		return err
	}

	for _, issue := range issuesToSend {
		err = appendIssue(issue)
		if err != nil {
			return err
		}
	}
	return nil
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

	patterns := strings.Split(os.Getenv("WATCH_PATTERNS"), ",")

	allIssues := []*github.Issue{}

	fmt.Println(patterns)

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

		allIssues = append(allIssues, issues...)
	}

	fmt.Printf("Found %d issues\n", len(allIssues))

	issuesToSend, err := filterSentIssues(allIssues)

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
