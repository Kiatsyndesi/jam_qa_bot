package jira_client

import (
	"errors"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"os"
)

func MakeNewJiraClient() (*JiraClient, error) {
	jiraCreds := new(JiraCredentials)

	jiraCreds.URL = os.Getenv("JIRA_URL")
	jiraCreds.login = os.Getenv("JIRA_LOGIN")
	jiraCreds.pass = os.Getenv("JIRA_PASS")

	tp := jira.BasicAuthTransport{
		Username: jiraCreds.login,
		Password: jiraCreds.pass,
	}

	client, err := jira.NewClient(tp.Client(), jiraCreds.URL)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("%v", err))
	}

	return &JiraClient{
		API: *client,
	}, nil
}
