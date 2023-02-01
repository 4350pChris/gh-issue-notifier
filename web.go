package main

import (
	"html/template"
	"log"
	"net/http"
)

type WebIssue struct {
	HtmlUrl string
	Title   string
	Number  int
	Updated string
}

type WebRepo struct {
	HtmlUrl      string
	FullName     string
	Description  string
	Issues       []WebIssue
	PatternLabel string
}

type WebData struct {
	Repos []WebRepo
}

func StartServer(config Config) error {
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	client := CreateClient(&config.GithubToken)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := WebData{
			Repos: []WebRepo{},
		}

		for _, pattern := range config.Patterns {
			repoRes, err := GetRepoForPattern(client, pattern)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			issues, err := GetIssuesForPattern(client, pattern)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			webIssues := []WebIssue{}

			for _, issue := range issues {
				webIssues = append(webIssues, WebIssue{
					HtmlUrl: issue.GetHTMLURL(),
					Title:   issue.GetTitle(),
					Number:  issue.GetNumber(),
					Updated: issue.GetUpdatedAt().Format("01/02/2006"),
				})
			}

			data.Repos = append(data.Repos, WebRepo{
				HtmlUrl:      repoRes.GetHTMLURL(),
				FullName:     repoRes.GetFullName(),
				Description:  repoRes.GetDescription(),
				Issues:       webIssues,
				PatternLabel: pattern.Label,
			})
		}

		tmpl.Execute(w, data)
	})

	log.Default().Println("Server listening on http://localhost:8080")
	return http.ListenAndServe(":8080", nil)
}
