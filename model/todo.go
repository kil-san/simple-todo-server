package model

type (
	Todo struct {
		Id     string `json:"id"`
		Title  string `json:"title"`
		Status string `json:"status"`
	}
)
