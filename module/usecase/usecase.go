package usecase

import (
	"encoding/csv"
	"fmt"
	"log"
	"ms-sv-jira/helper/logger"
	"ms-sv-jira/models/db"
	"ms-sv-jira/module"
	"os"
	"reflect"
	"strconv"
	"time"
)

type Usecase struct {
	repo    module.Repository
	timeout time.Duration
	log     logger.Logger
}

func NewUsecase(repo module.Repository, timeout time.Duration, log logger.Logger) module.Usecase {
	return &Usecase{
		repo:    repo,
		timeout: timeout,
		log:     log,
	}
}

func (usecase *Usecase) GetIssuesByProjectId(id string) (res []db.Task, err error) {
	data, err := usecase.repo.GetIssuesByProjectId(id)
	if err != nil {
		return nil, err
	}

	exists := make(map[string]bool)
	for _, v := range data {
		if v.Type == "Task" && !exists[v.JiraIssueKey] {
			res = append(res, db.Task{
				Id:        v.JiraIssueKey,
				ProjectId: v.ProjectId,
				BoardId:   v.BoardId,
				SprintId:  v.SprintId,
				Status: db.Status{
					Name:      v.Status,
					UpdatedAt: v.UpdatedAt,
				},
				Summary:     v.Summary,
				DisplayName: v.DisplayName,
				Priority:    v.Priority,
				CreatedAt:   v.CreatedAt,
			})
			exists[v.JiraIssueKey] = true
		}
	}

	for i, t := range res { // task
		exists := make(map[string]bool)
		for _, d := range data { // subtask
			if d.Type == "Subtask" && d.ParentKey == t.Id && !exists[d.JiraIssueKey] {
				res[i].SubTask = append(res[i].SubTask, db.Issue{
					JiraIssueKey: d.JiraIssueKey,
					Summary:      d.Summary,
					Status: db.Status{
						Name:      d.Status,
						UpdatedAt: d.UpdatedAt,
					},
					DisplayName: d.DisplayName,
					Priority:    d.Priority,
					CreatedAt:   d.CreatedAt,
				})
				exists[d.JiraIssueKey] = true
			}
		}
	}

	return res, err
}

func (u *Usecase) Csv(table string) {
	var (
		data     []map[string]string
		header   []string
		filename string
	)
	switch table {
	case "projects":
		projects, err := u.repo.GetProjects()
		if err != nil {
			return
		}

		for _, v := range projects {
			value := map[string]string{
				"id":           strconv.Itoa(v.ProjectID),
				"project_key":  v.JiraProjectKey,
				"name":         v.Name,
				"description":  v.Description,
				"lead":         v.Lead,
				"type_project": v.TypeProject,
				"offload_time": v.CreatedAt.Format("2006-01-02"),
				"created_at":   v.Created.Format("2006-01-02"),
			}
			data = append(data, value)
		}

		header = []string{"id",
			"project_key",
			"name",
			"description",
			"lead",
			"type_project",
			"offload_time",
			"created_at"}
		filename = "projects.csv"
	case "issues":
		issues, err := u.repo.GetIssues()
		if err != nil {
			return
		}

		for _, v := range issues {
			value := map[string]string{
				"issue_id":       strconv.Itoa(v.IssueID),
				"jira_issue_key": v.JiraIssueKey,
				"project_id":     v.ProjectID,
				"epic_key":       v.EpicKey,
				"summary":        v.Summary,
				"description":    v.Description,
				"type":           v.Type,
				"status":         v.Status,
				"assignee_id":    v.AssigneeID,
				"reporter_id":    v.ReporterID,
				"priority":       v.Priority,
				"created_at":     v.CreatedAt.Format("2006-01-02"),
				"updated_at":     v.UpdatedAt.Format("2006-01-02"),
				"created":        v.Created.Format("2006-01-02"),
				"updated":        v.Updated.Format("2006-01-02"),
				"parent_key":     v.ParentKey,
			}
			data = append(data, value)
		}

		header = []string{"issue_id", "jira_issue_key", "project_id", "epic_key", "summary", "description", "type", "status", "assignee_id", "reporter_id", "priority", "created_at", "updated_at", "created", "updated", "parent_key"}
		filename = "issues.csv"
	case "users":
		users, err := u.repo.GetUsers()
		if err != nil {
			return
		}

		for _, v := range users {
			active := "false"
			if v.Active {
				active = "true"
			}
			data = append(data, map[string]string{
				"user_id":      strconv.Itoa(v.UserID),
				"jira_user_id": v.JiraUserID,
				"display_name": v.DisplayName,
				"email":        v.Email,
				"active":       active,
				"created":      v.Created.Format("2006-01-02"),
			})
		}

		header = []string{"user_id", "jira_user_id", "display_name", "email", "active", "created"}
		filename = "users.csv"
	}

	exportToCSV(data, header, filename)
}

func extractJSONTags(model interface{}) []string {
	var headers []string
	t := reflect.TypeOf(model).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			headers = append(headers, jsonTag)
		}
	}
	return headers
}

// Helper function to extract field values dynamically
func extractFieldValues(record interface{}) []string {
	var values []string
	v := reflect.ValueOf(record).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		values = append(values, fmt.Sprintf("%v", field.Interface()))
	}
	return values
}

func exportToCSV(records []map[string]string, headers []string, fileName string) {
	// Create a CSV file
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("failed to create csv file: %v", err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(headers); err != nil {
		log.Fatalf("failed to write headers: %v", err)
	}

	// Query data using the provided model type
	// modelType := reflect.New(reflect.TypeOf(records).Elem()).Interface()

	// Write data to CSV dynamically
	for _, record := range records {
		var row []string
		for _, v := range headers {
			if val, ok := record[v]; ok {
				row = append(row, val)
			}
		}
		// row := extractFieldValues(&record)
		if err := writer.Write(row); err != nil {
			log.Fatalf("failed to write record to csv: %v", err)
		}
	}

	fmt.Println("Data exported to CSV successfully")
}
