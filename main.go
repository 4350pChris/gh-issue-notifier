package main

import (
	"log"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	config := GetConfig()
	for {
		err := doWork(config)
		if err != nil {
			log.Fatal(err)
		}

		log.Default().Printf("Sleeping for %v\n", config.Interval)
		time.Sleep(config.Interval)
	}
}

func doWork(config Config) error {
	repos, err := GetReposWithIssues(config.Patterns)

	if err != nil {
		return err
	}

	if len(repos) == 0 {
		log.Default().Println("No issues to send")
		return nil
	}

	for _, repo := range repos {
		issuesToSend, err := filterSentIssues(repo)
		if err != nil {
			return err
		}
		// replace issues with only the ones that need to be sent
		repo.issues = issuesToSend
	}
	err = NotifyOfIssues(config.ShoutrrrUrl, repos)

	if err != nil {
		return err
	}

	for _, repo := range repos {
		err = StoreIssues(repo)
		if err != nil {
			log.Default().Printf("Error storing issues for repo %v: %v", repo.fullName, err)
			return err
		}
	}
	return nil
}
