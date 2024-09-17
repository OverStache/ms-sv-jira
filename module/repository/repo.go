package repository

import (
	"ms-sv-jira/helper/logger"
	"ms-sv-jira/models/db"
	"ms-sv-jira/module"

	"gorm.io/gorm"
)

type Repository struct {
	Conn *gorm.DB
	log  logger.Logger
}

func NewRepository(Conn *gorm.DB, log logger.Logger) module.Repository {
	return &Repository{Conn, log}
}

func (repository *Repository) GetIssuesByProjectId(id string) (res []db.SelectIssue, err error) {
	err = repository.Conn.Raw(`SELECT b.jira_board_id , s.jira_sprint_id , i.*, u.display_name
                        FROM issues i
                        JOIN issue_sprint_link isl ON i.jira_issue_key = isl.issue_id
                        JOIN sprints s ON isl.sprint_id = s.jira_sprint_id 
                        JOIN boards b ON s.board_id = b.jira_board_id 
                        left join users u on i.assignee_id = u.jira_user_id 
                        WHERE b.project_id = ?`, id).
		Find(&res).Error
	return
}

func (repository *Repository) GetProjects() (res []db.Projects, err error) {
	err = repository.Conn.Find(&res).Error
	return
}

func (repository *Repository) GetIssues() (res []db.Issues, err error) {
	err = repository.Conn.Find(&res).Error
	return
}

func (repository *Repository) GetUsers() (res []db.Users, err error) {
	err = repository.Conn.Find(&res).Error
	return
}

func (repository *Repository) GetAttachments() (res []db.Attachments, err error) {
	err = repository.Conn.Find(&res).Error
	return
}

func (repository *Repository) GetComments() (res []db.Comments, err error) {
	err = repository.Conn.Find(&res).Error
	return
}

func (repository *Repository) GetBoards() (res []db.Boards, err error) {
	err = repository.Conn.Find(&res).Error
	return
}

func (repository *Repository) GetSprints() (res []db.Sprints, err error) {
	err = repository.Conn.Find(&res).Error
	return
}

func (repository *Repository) GetIssueLinks() (res []db.IssueLinks, err error) {
	err = repository.Conn.Find(&res).Error
	return
}
