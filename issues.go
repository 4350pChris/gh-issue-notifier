package main

import (
	"context"
	"log"

	"github.com/google/go-github/v50/github"
)

func GetReposWithIssues(patterns []WatchPattern) ([]*Repository, error) {
	client := github.NewClient(nil)

	repos := []*Repository{}

	for _, pattern := range patterns {

		repo, _, err := client.Repositories.Get(context.Background(), pattern.Owner, pattern.Repo)
		if err != nil {
			return nil, err
		}

		issues, _, err := client.Issues.ListByRepo(context.Background(), pattern.Owner, pattern.Repo, &github.IssueListByRepoOptions{Labels: []string{pattern.Label}})
		if err != nil {
			return nil, err
		}

		dbRepo := ghIssuesToDbRepo(repo, issues)

		issuesToSend, err := filterSentIssues(dbRepo)
		if err != nil {
			return nil, err
		}

		log.Default().Printf("Found %d issues for %s, out of which %d are new", len(issues), pattern.Repo, len(issuesToSend))

		// skip repos that have no fitting issues
		if len(issues) > 0 {
			repos = append(repos, dbRepo)
		}
	}

	return repos, nil
}

func ghIssuesToDbRepo(repo *github.Repository, issues []*github.Issue) *Repository {
	convertedIssues := []Issue{}

	for _, issue := range issues {
		convertedIssues = append(convertedIssues, Issue{
			id:     issue.GetID(),
			title:  issue.GetTitle(),
			number: issue.GetNumber(),
			repoId: issue.GetRepository().GetID(),
		})
	}

	return &Repository{
		id:          repo.GetID(),
		fullName:    repo.GetFullName(),
		description: repo.GetDescription(),
		htmlUrl:     repo.GetHTMLURL(),
		issues:      convertedIssues,
	}
}

func filterSentIssues(repo *Repository) ([]Issue, error) {
	issuesToSend := []Issue{}

	for _, issue := range repo.issues {
		found, err := IssueExists(issue)
		if err != nil {
			return nil, err
		}

		if !found {
			issuesToSend = append(issuesToSend, issue)
		}
	}
	return issuesToSend, nil
}
