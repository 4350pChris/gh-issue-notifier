package main

import (
	"log"
	"net/http"
	"text/template"
)

type WebData struct {
	Repos []Repository
}

func StartServer() error {
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		repos, err := RetrieveReposWithIssues()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, WebData{Repos: repos})
	})

	log.Default().Println("Server listening on http://localhost:8080")
	return http.ListenAndServe(":8080", nil)
}
