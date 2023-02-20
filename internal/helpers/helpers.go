package helpers

import (
	"github.com/andygrunwald/go-jira"
	"github.com/mattermost/mattermost-server/v5/model"
	"jam_qa_bot/internal/app/issue-validators"
	"jam_qa_bot/internal/app/jira-client"
	"log"
	"os"
)

// FindIssueIDsWithPoorDescription - function for finding issue keys where description less const from issue validators
func FindIssueIDsWithPoorDescription(client *jira_client.JiraClient) (map[string]string, *jira.Response, error) {
	jql := "project = " + os.Getenv("JIRA_PROJECT_NAME") + "AND resolution = Unresolved"

	issues, resp, err := client.API.Issue.Search(jql, nil)

	if err != nil {
		log.Fatal(err)
	}

	issuesIDs := make(map[string]string)

	for _, issue := range issues {
		if !issue_validators.IsDescriptionHasNormalLength(&issue) {
			issuesIDs[issue.Key] = issue.Fields.Creator.Name
		}
	}

	return issuesIDs, resp, nil
}

// PrintError custom error for mattermost client model
func PrintError(err *model.AppError) {
	println("\tError Details:")
	println("\t\t" + err.Message)
	println("\t\t" + err.Id)
	println("\t\t" + err.DetailedError)
}
