package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

var (
	location, _ = time.LoadLocation("Europe/Moscow")
)

// RFC 9457
type ProblemDetails struct {
	StatusCode int    `json:"status_code" example:"400"`
	Method     string `json:"method" example:"GET"`
	Time       string `json:"time" example:"2006-01-02 15:04:05"`
	Type       string `json:"type" example:"http://api/v1/problems"`
	Title      string `json:"title" example:"Invalid input"`
	Details    string `json:"details" example:"field 'email' is invalid"`
	Instance   string `json:"instance" example:"api/v1/register"`
}

func HttpRespErrRFC9457(handler, title string, err error, statusCode int, w http.ResponseWriter, r *http.Request, logger *CustomLogger) {
	logger.Error("Handler: %s | Title: %s | Error: %s | Path: %s | Method: %s", handler, title, err.Error(), r.URL.Path, r.Method)

	problem := ProblemDetails{
		StatusCode: statusCode,
		Method:     r.Method,
		Time:       time.Now().In(location).Format("2006-01-02 15:04:05"),
		Type:       os.Getenv("ERR_INDENT_ADDR"),
		Title:      title,
		Details:    err.Error(),
		Instance:   r.Pattern,
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(problem)
}
