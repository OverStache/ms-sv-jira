package module

import (
	"ms-sv-jira/models/db"
)

type Repository interface {
	GetIssuesByProjectId(id string) ([]db.SelectIssue, error)
	GetProjects() (res []db.Projects, err error)
	GetIssues() (res []db.Issues, err error)
	GetUsers() (res []db.Users, err error)
}
