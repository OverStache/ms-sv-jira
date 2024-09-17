package module

import (
	"ms-sv-jira/models/db"
)

type Usecase interface {
	GetIssuesByProjectId(id string) ([]db.Task, error)
	Csv(table string) error
	GetProjects() (res []db.Projects, err error)
	GetIssues() (res []db.Issues, err error)
	GetUsers() (res []db.Users, err error)
	GetAttachments() (res []db.Attachments, err error)
	GetComments() (res []db.Comments, err error)
	GetBoards() (res []db.Boards, err error)
	GetSprints() (res []db.Sprints, err error)
	GetIssueLinks() (res []db.IssueLinks, err error)
}
