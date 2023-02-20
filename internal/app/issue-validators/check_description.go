package issue_validators

import (
	"github.com/andygrunwald/go-jira"
)

const validDescriptionLength = 300

func IsDescriptionHasNormalLength(issue *jira.Issue) bool {
	issueLength := len(issue.Fields.Description)

	if issueLength <= validDescriptionLength {
		return false
	}

	return true
}
