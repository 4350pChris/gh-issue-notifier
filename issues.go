package main

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

func CreateClient(token *string) *github.Client {
	if token == nil {
		return github.NewClient(nil)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func GetIssuesForPattern(client *github.Client, pattern WatchPattern) ([]*github.Issue, error) {
	issues, _, err := client.Issues.ListByRepo(context.Background(), pattern.Owner, pattern.Repo, &github.IssueListByRepoOptions{Labels: []string{pattern.Label}})
	if err != nil {
		return nil, err
	}

	return issues, nil
}

func GetIssuesToSend(client *github.Client, patterns []WatchPattern) ([]Sendable, error) {
	sendables := []Sendable{}

	for _, pattern := range patterns {

		issues, err := GetIssuesForPattern(client, pattern)
		if err != nil {
			return nil, err
		}

		issuesToSend, err := filterSentIssues(issues)
		if err != nil {
			return nil, err
		}

		log.Default().Printf("Found %d issues for %s, out of which %d are new", len(issues), pattern.Repo, len(issuesToSend))

		// skip repos that have no new issues
		if len(issuesToSend) > 0 {
			sendables = append(sendables, Sendable{repo: pattern.Repo, issues: issuesToSend})
		}
	}

	return sendables, nil
}

func filterSentIssues(issues []*github.Issue) ([]*github.Issue, error) {
	content, err := os.ReadFile("issues.txt")

	if err != nil && errors.Is(err, fs.ErrNotExist) {
		log.Default().Println("Creating issues.txt file... ")
		err = os.WriteFile("issues.txt", []byte(""), 0644)
	}

	if err != nil {
		return nil, err
	}

	readIds := strings.Split(string(content), "\n")

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

func RememberIssue(issue *github.Issue) error {
	content, err := os.ReadFile("issues.txt")

	if err != nil {
		return err
	}

	content = append(content, []byte(*issue.URL+"\n")...)

	err = os.WriteFile("issues.txt", content, 0644)

	if err != nil {
		return err
	}

	return nil
}
