package main

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v50/github"
)

func GetIssuesToSend(patterns []WatchPattern) ([]Sendable, error) {
	client := github.NewClient(nil)

	sendables := []Sendable{}

	for _, pattern := range patterns {

		issues, _, err := client.Issues.ListByRepo(context.Background(), pattern.Owner, pattern.Repo, &github.IssueListByRepoOptions{Labels: []string{pattern.Label}})
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

func RememberIssue(issue *github.Issue) error {
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
