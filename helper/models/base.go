package models

import (
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	ApiHandlerFailedMethod      = "Invalid Method Middleware"
	ApiHandlerFailedToken       = "Invalid Token Middleware"
	ApiHandlerFailedContentType = "Invalid Content Type Middleware"
	ApiHandlerFailedRoute       = "Invalid Route Middleware"
)

type Count struct {
	Count int `json:"count"`
}
type ResponseCommand struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}
type ResponseCommandResult struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}
type MetaPagination struct {
	Page          int `json:"page"`
	Total         int `json:"total_pages"`
	TotalRecords  int `json:"total_records"`
	Prev          int `json:"prev"`
	Next          int `json:"next"`
	RecordPerPage int `json:"record_per_page"`
}

type Response struct {
	StatusCode int         `json:"status_code"`
	Status     string      `json:"status_desc"`
	Msg        string      `json:"message"`
	Data       interface{} `json:"data"`
	Errors     interface{} `json:"errors"`
	RequestId  string      `json:"request_id"`
}

type Paginator struct {
	CurrentPage  int `json:"current_page"`
	PerPage      int `json:"records_per_page"`
	PreviousPage int `json:"back_page"`
	NextPage     int `json:"next_page"`
	TotalRecords int `json:"total_records"`
	TotalPages   int `json:"total_pages"`
}

func (r *Response) MappingResponseSuccess(message string, data interface{}, ctx echo.Context, pernr string, role string, apiHandler string) {
	r.StatusCode = http.StatusOK
	r.Status = http.StatusText(r.StatusCode)
	r.Msg = message
	r.Data = data
	r.Errors = nil
	r.RequestId = ctx.Get("REQUEST-ID").(string)

	ctx.Set("RESPONSE", r)
	ctx.Set("API_HANDLER", apiHandler)
	ctx.Set("ERROR", "")
	ctx.Set("PERNR", pernr)
	ctx.Set("ROLE", role)
}

func (r *Response) MappingResponseError(statusCode int, message string, error interface{}, ctx echo.Context, pernr string, role string, apiHandler string) {
	r.StatusCode = statusCode
	r.Status = http.StatusText(r.StatusCode)
	r.Msg = message
	r.Data = nil
	r.Errors = error
	r.RequestId = ctx.Get("REQUEST-ID").(string)

	ctx.Set("RESPONSE", r)
	ctx.Set("API_HANDLER", apiHandler)
	if r.Errors == nil {
		r.Errors = ""
	}
	ctx.Set("ERROR", r.Errors)
	ctx.Set("PERNR", pernr)
	ctx.Set("ROLE", role)
}
func (p Paginator) MappingPaginator(page, limit, perPageRecords, totalAllRecords int) Paginator {
	totalPage := int(math.Ceil(float64(totalAllRecords) / float64(limit)))
	prev := page
	next := page
	if page != 1 {
		prev = page - 1
	}

	if page != totalPage {
		next = page + 1
	}

	p = Paginator{
		CurrentPage:  page,
		PerPage:      perPageRecords,
		PreviousPage: prev,
		NextPage:     next,
		TotalRecords: totalAllRecords,
		TotalPages:   totalPage,
	}

	return p
}
