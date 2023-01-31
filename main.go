package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/go-github/v50/github"
	_ "github.com/joho/godotenv/autoload"
)

type Sendable struct {
	repo   string
	issues []*github.Issue
}

func main() {
	config := GetConfig()
	for {
		err := doWork(config)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Sleeping for %v\n", config.Interval)
		time.Sleep(config.Interval)
	}
}

func doWork(config Config) error {
	sendables, err := GetIssuesToSend(config.Patterns)

	if err != nil {
		return err
	}

	if len(sendables) == 0 {
		fmt.Println("No issues to send")
		return nil
	}

	err = NotifyOfIssues(config.ShoutrrrUrl, sendables)

	if err != nil {
		return err
	}

	// save all issues to txt file
	for _, sendable := range sendables {
		for _, issue := range sendable.issues {
			err = RememberIssue(issue)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
