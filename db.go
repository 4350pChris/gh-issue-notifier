package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Issue struct {
	id      int64
	title   string
	number  int
	repoId  int64
	htmlUrl string
}

type Repository struct {
	id          int64
	fullName    string
	description string
	htmlUrl     string
	issues      []Issue
}

func getDbConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableStmt := `create table if not exists repositories (id integer not null primary key, full_name text not null, description text, html_url text);
	create table if not exists issues (id integer not null primary key, title text, number integer, repository_id integer, foreign key (repository_id) references repositories(id) );`

	_, err = db.Exec(createTableStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, createTableStmt)
		return nil, err
	}

	return db, nil
}

func RetrieveIssues() ([]Issue, error) {
	db, err := getDbConnection()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("select * from issues")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	issues := []Issue{}
	for rows.Next() {
		var issue Issue
		err = rows.Scan(&issue.id, &issue.title, &issue.number, &issue.repoId)
		if err != nil {
			return nil, err
		}
		issues = append(issues, issue)
	}
	return issues, nil
}

func StoreIssues(repo *Repository) error {
	db, err := getDbConnection()
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	repoStmt, err := tx.Prepare("insert into repositories (id, full_name, description, html_url) values (?, ?, ?, ?) on conflict(id) do update set full_name = excluded.full_name, description = excluded.description, html_url = excluded.html_url")
	if err != nil {
		return err
	}
	defer repoStmt.Close()

	issueStmt, err := tx.Prepare("insert into issues (id, repository_id, title, number) values (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer issueStmt.Close()

	for _, issue := range repo.issues {
		_, err := repoStmt.Exec(repo.id, repo.fullName, repo.description, repo.htmlUrl)
		if err != nil {
			return err
		}
		_, err = issueStmt.Exec(issue.id, repo.id, issue.title, issue.number)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func IssueExists(issue Issue) (bool, error) {
	db, err := getDbConnection()
	if err != nil {
		return false, err
	}
	rows, err := db.Query("select * from issues where id = ?", issue.id)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}
