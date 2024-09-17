package module

import (
	"ms-sv-jira/models/db"
)

type Usecase interface {
	GetIssuesByProjectId(id string) ([]db.Task, error)
	Csv(table string)
}
