package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"ms-sv-jira/helper"
	"ms-sv-jira/helper/logger"
	helperModels "ms-sv-jira/helper/models"
	"ms-sv-jira/models"
	"net/http"
	"net/url"
	"os"
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

			go m.AuditTrails(context.Background(), userAuditTrail)
		}

		l.Info("[" + strconv.Itoa(c.Response().Status) + "] " + "[" + c.Request().Method + "] " + c.Request().Host + c.Request().URL.String())

		l.Info("Closing")
		return nil
	}
}

// CORSValidationGlobalResponse will handle the CORSValidationGlobalResponse middleware
func (m *GoMiddleware) CORSValidationGlobalResponse(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Method == "POST" &&
			c.Request().Header.Get(echo.HeaderContentType) == echo.MIMEApplicationJSON {
			c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			response := new(helperModels.Response)
			response.MappingResponseError(http.StatusBadRequest, "Invalid "+
				echo.HeaderContentType+" "+echo.MIMEApplicationJSON, nil, c, "", "", helperModels.ApiHandlerFailedContentType)
			c.JSON(response.StatusCode, response)
			return nil
		}

		if c.Request().Method == "POST" &&
			c.Request().Header.Get(echo.HeaderContentType) == "application-bribrain/json" {
			c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		}

		err := next(c)

		if err != nil {
			allowedMethod := []string{"GET", "POST", "DELETE"}
			if !helper.InArray(c.Request().Method, allowedMethod) {
				response := new(helperModels.Response)
				response.MappingResponseError(http.StatusBadRequest, "Invalid Method", nil, c, "", "", helperModels.ApiHandlerFailedMethod)
				c.JSON(response.StatusCode, response)
				return err
			} else if err.Error() == "code=404, message=Not Found" {
				response := new(helperModels.Response)
				response.MappingResponseError(http.StatusNotFound, "Routes Not Found", nil, c, "", "", helperModels.ApiHandlerFailedRoute)

				c.JSON(response.StatusCode, response)
				return err
			} else if err.Error() == "code=405, message=Method Not Allowed" {
				response := new(helperModels.Response)
				response.MappingResponseError(http.StatusMethodNotAllowed, "Method Not Allowed", nil, c, "", "", helperModels.ApiHandlerFailedMethod)

				c.JSON(response.StatusCode, response)
				return err
			}
		}
		return nil
	}
}

// Authentication will handle the Authentication Token middleware
func (m *GoMiddleware) Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		response := new(helperModels.Response)
		if c.Request().Header["Authorization"] != nil {
			roles := []string{"PAB", "KC", "KW", "KP"}
			accessLevel := []string{"Bripa"}
			//accessLevel := []string{}
			token := helper.StringToStringNullable(c.Request().Header.Get("Authorization"))
			resp, err := m.AuthValidate(c.Request().Context(), token, roles, accessLevel)
			if err != nil {
				response.MappingResponseError(http.StatusUnauthorized, err.Error(), nil, c, "", "", helperModels.ApiHandlerFailedToken)
				c.JSON(response.StatusCode, response)
				return err
			}
			if resp.StatusCode != http.StatusOK {
				response.MappingResponseError(resp.StatusCode, resp.Msg, resp.Errors, c, "", "", helperModels.ApiHandlerFailedToken)
				c.JSON(resp.StatusCode, resp)
				return nil
			}
			convertClaims := helper.ObjectToString(resp.Data)
			var jwtClaims *models.BRIBrainClaims
			if errUnmarshal := json.Unmarshal([]byte(string(convertClaims)), &jwtClaims); errUnmarshal != nil {
				response.MappingResponseError(http.StatusUnauthorized, helperModels.ErrUnAuthorize.Error(), nil, c, "", "", helperModels.ApiHandlerFailedToken)
				c.JSON(response.StatusCode, response)

				return helperModels.ErrUnAuthorize
			}
			c.Set("JWT_CLAIMS", jwtClaims)
			generalRequest := jwtClaims.MappingToGeneralRequestV3()
			c.Set("GENERAL_REQUEST_V3", generalRequest)
			next(c)
		} else {
			response.MappingResponseError(http.StatusUnauthorized, helperModels.ErrUnAuthorize.Error(), nil, c, "", "", helperModels.ApiHandlerFailedToken)
			c.JSON(response.StatusCode, response)

			return helperModels.ErrUnAuthorize
		}
		return nil
	}
}

// Authentication will handle the Authentication Token middleware
func (m *GoMiddleware) AuthValidate(ctx context.Context, token *string, roles []string, accessLevel []string) (response *helperModels.Response, err error) {
	l := logger.L
	validateAccess := models.ValidateAccess{
		Roles:       roles,
		AccessLevel: accessLevel,
	}
	b, err := json.Marshal(validateAccess)
	if err != nil {
		l.Error("GoMiddleware.AuthValidate: %s", err.Error())
		return nil, helperModels.ErrUnAuthorize
	}

	url, err := url.Parse(os.Getenv("VALIDATE_TOKEN_URL"))
	if err != nil {
		l.Error("GoMiddleware.AuthValidate: %s", err.Error())
		return nil, helperModels.ErrUnAuthorize
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url.String(), bytes.NewReader(b))
	if err != nil {
		l.Error("GoMiddleware.AuthValidate: %s", err.Error())
		return nil, helperModels.ErrUnAuthorize
	}

	req.Header.Add("Content-Type", "application-bribrain/json")
	req.Header.Add("Authorization", helper.StringNullableToString(token))

	// tr := &http.Transport{
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// }
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		l.Error("GoMiddleware.AuthValidate: %s", err.Error())
		return nil, helperModels.ErrUnAuthorize
	}
	defer resp.Body.Close()

	br := &helperModels.Response{}
	if err := json.NewDecoder(resp.Body).Decode(br); err != nil {
		l.Error("GoMiddleware.AuthValidate: %s", err.Error())
		return nil, helperModels.ErrUnAuthorize
	}

	return br, nil

}

// AuditTrails will handle the AuditTrails Token middleware
func (m *GoMiddleware) AuditTrails(ctx context.Context, body models.BribrainAuditTrailRequest) (response *helperModels.Response, err error) {
	l := logger.L

	b, err := json.Marshal(body)
	if err != nil {
		l.Error("GoMiddleware.AuditTrails: %s", err.Error())
		return nil, helperModels.ErrUnAuthorize
	}

	url, err := url.Parse(os.Getenv("AUDIT_TRAILS_URL"))
	if err != nil {
		l.Error("GoMiddleware.AuditTrails: %s", err.Error())
		return nil, helperModels.ErrUnAuthorize
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url.String(), bytes.NewReader(b))
	if err != nil {
		l.Error("GoMiddleware.AuditTrails: %s", err.Error())
		return nil, helperModels.ErrUnAuthorize
	}

	req.Header.Add("Content-Type", "application-bribrain/json")

	// tr := &http.Transport{
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// }
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		l.Error("GoMiddleware.AuditTrails: %s", err.Error())
		return nil, helperModels.ErrUnAuthorize
	}
	defer resp.Body.Close()

	br := &helperModels.Response{}
	if err := json.NewDecoder(resp.Body).Decode(br); err != nil {
		l.Error("GoMiddleware.AuditTrails: %s", err.Error())
		return nil, helperModels.ErrUnAuthorize
	}

	return br, nil

}

// InitMiddleware intialize the middleware
func InitMiddleware() *GoMiddleware {
	return &GoMiddleware{}
}
