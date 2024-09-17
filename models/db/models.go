package db

import "time"

//from db

type SelectIssue struct { // select
	ID           string `gorm:"column:issue_id"`
	JiraIssueKey string
	ProjectId    string
	BoardId      string `gorm:"column:jira_board_id"`
	SprintId     string `gorm:"column:jira_sprint_id"`
	Summary      string
	Status       string
	DisplayName  string
	Type         string
	Priority     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ParentKey    string

	// Child     []Child `json:"child,omitempty"` // Optional, only present for some issues
}

type Issue struct { // select
	JiraIssueKey string    `json:"id"`
	Summary      string    `json:"summary"`
	Status       Status    `json:"status"`
	DisplayName  string    `json:"asignee"`
	Priority     string    `json:"priority"`
	CreatedAt    time.Time `json:"created_at"`

	// Child     []Child `json:"child,omitempty"` // Optional, only present for some issues
}

type Task struct {
	Id          string    `json:"id"`
	ProjectId   string    `json:"project_id"`
	BoardId     string    `json:"board_id"`
	SprintId    string    `json:"sprint_id"`
	Summary     string    `json:"summary"`
	Status      Status    `json:"status"`
	DisplayName string    `json:"asignee"`
	Priority    string    `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	SubTask     []Issue   `json:"sub_tasks"`
}

type Status struct {
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Projects struct {
	ProjectID      int
	JiraProjectKey string
	Name           string
	Description    string
	Lead           string
	TypeProject    string
	CreatedAt      time.Time
	Created        time.Time
}

type Users struct {
	UserID      int
	JiraUserID  string
	DisplayName string
	Email       string
	Active      bool
	Created     time.Time
}

type Issues struct {
	IssueID      int
	JiraIssueKey string
	ProjectID    string
	EpicKey      string
	Summary      string
	Description  string
	Type         string
	Status       string
	AssigneeID   string
	ReporterID   string
	Priority     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Created      time.Time
	Updated      time.Time
	ParentKey    string
}

type Sprints struct {
	SprintID     int
	JiraSprintID string
	BoardID      int
	Name         string
	Goal         string
	State        string
	StartDate    time.Time
	EndDate      time.Time
	CompleteDate time.Time
	CreatedAt    time.Time
	Created      time.Time
}

type Attachments struct {
	AttachmentID     int
	IssueID          string
	FileName         string
	MimeType         string
	FileSize         int
	JiraAttachmentID string
	CreatedAt        time.Time
	Created          time.Time
}

type Boards struct {
	BoardID     int
	JiraBoardID string
	ProjectID   string
	Name        string
	Type        string
	CreatedAt   time.Time
	Created     time.Time
}

type Comments struct {
	CommentID int
	IssueID   string
	AuthorID  string
	Body      string
	CreatedAt time.Time
	Created   time.Time
}

type IssueLinks struct {
	LinkID         int
	IssueID        string
	LinkedIssueKey string
	URL            string
	CreatedAt      time.Time
	Title          string
}
