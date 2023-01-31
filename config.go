package main

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type WatchPattern struct {
	Owner string `yaml:"owner"`
	Repo  string `yaml:"repo"`
	Label string `yaml:"label"`
}

type PatternFile struct {
	Patterns []WatchPattern `yaml:"patterns"`
}

type Config struct {
	Interval    time.Duration
	Patterns    []WatchPattern
	ShoutrrrUrl string
}

func parsePatterns() []WatchPattern {
	content, err := os.ReadFile("patterns.yaml")
	if err != nil {
		log.Fatal(err)
	}

	parsed := PatternFile{}

	err = yaml.Unmarshal(content, &parsed)

	if err != nil {
		log.Fatal(err)
	}

	return parsed.Patterns
}

func parseInterval() time.Duration {
	duration, err := time.ParseDuration(os.Getenv("INTERVAL"))
	if err != nil {
		log.Fatal(err)
	}
	return duration
}

func GetConfig() Config {
	patterns := parsePatterns()
	interval := parseInterval()

	return Config{
		Interval:    interval,
		Patterns:    patterns,
		ShoutrrrUrl: os.Getenv("SHOUTRRR_URL"),
	}
}
