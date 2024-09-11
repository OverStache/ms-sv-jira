package delivery_jira

import "github.com/gin-gonic/gin"

type JiraDelivery interface {
	GetAllProjectDelivery(ginContext *gin.Context)
	GetAllSprintDelivery(ginContext *gin.Context)
	GetAllCardDelivery(ginContext *gin.Context)
}
