package middleware

import (
	"encoding/json"
	"ms-sv-jira/helper"
	"ms-sv-jira/helper/logger"
	"ms-sv-jira/models"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// GoMiddleware represent the data-struct for middleware
type GoMiddleware struct {
	// another stuff , may be needed by middleware
}

// CORS will handle the CORS middleware
func (m *GoMiddleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Request().Header.Set("Access-Control-Allow-Headers", "*")
		c.Request().Header.Set("Access-Control-Allow-Origin", "*")
		c.Request().Header.Set("Access-Control-Allow-Methods", "*")

		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set("Access-Control-Allow-Headers", "*")
		c.Response().Header().Set("Access-Control-Allow-Methods", "*")
		return next(c)
	}
}

// LOG will handle the LOG middleware
func (m *GoMiddleware) Log(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		l := logger.L
		l.Info("Accepted")

		requestId := uuid.New().String()
		c.Set("REQUEST-ID", requestId)

		userAuditTrail := models.BribrainAuditTrailRequest{}
		payloadAuditTrail := models.PayloadBribrainAuditTrail{}
		headersAuditTrail := models.HeadersPayloadBribrainAuditTrail{}

		next(c)

		var resp interface{}
		var apiHandler string
		var pernr string
		var role string
		if c.Get("RESPONSE") != nil {
			resp = c.Get("RESPONSE").(interface{})
		}
		if c.Get("API_HANDLER") != nil {
			apiHandler = c.Get("API_HANDLER").(string)
		}
		if c.Get("PERNR") != nil {
			pernr = c.Get("PERNR").(string)
		}
		if c.Get("ROLE") != nil {
			role = c.Get("ROLE").(string)
		}
		//err := c.Get("ERROR").(string)
		statusCode := c.Response().Status
		reqbody, _ := json.Marshal(c.Request().Body)

		if apiHandler != "" {
			remarkAuditTrail := helper.IntToString(statusCode)
			headersAuditTrail = headersAuditTrail.MappingHeadersPayloadBribrainAuditTrail(c.Request().Header.Get(echo.HeaderContentType), c.Request().Header.Get("Authorization"))
			payloadAuditTrail = payloadAuditTrail.MappingPayloadBribrainAuditTrail(c.Request().Host+c.Request().URL.Path, headersAuditTrail, string(reqbody), helper.JsonString(resp), statusCode)
			userAuditTrail = userAuditTrail.MappingBribrainAuditTrail(apiHandler, remarkAuditTrail, pernr, role, helper.JsonString(payloadAuditTrail))
		}

		l.Info("[" + strconv.Itoa(c.Response().Status) + "] " + "[" + c.Request().Method + "] " + c.Request().Host + c.Request().URL.String())

		l.Info("Closing")
		return nil
	}
}

// InitMiddleware intialize the middleware
func InitMiddleware() *GoMiddleware {
	return &GoMiddleware{}
}
