package jira_client

import "github.com/andygrunwald/go-jira"

type JiraCredentials struct {
	URL   string
	login string
	pass  string
}

type JiraClient struct {
	API jira.Client
}
