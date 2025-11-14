package models

import (
	"net/http"
	"strconv"
)

type PaginatedQuery struct {
	Limit  int `json:"limit,omitempty" validate:"min=1,max=100"`
	Offset int `json:"offset,omitempty" validate:"min=0"`
}

func (pq *PaginatedQuery) Parse(r *http.Request) error {
	q := r.URL.Query()

	if limit := q.Get("limit"); limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return err
		}
		pq.Limit = l
	}

	if offset := q.Get("offset"); offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return err
		}
		pq.Offset = o
	}
	return nil
}

func (pq *PaginatedQuery) SetDefaults() {
	if pq.Limit == 0 {
		pq.Limit = 20
	}
	if pq.Offset == 0 {
		pq.Offset = 0
	}
}

type JobFilter struct {
	Position       string   `json:"position,omitempty"`
	Country        string   `json:"country,omitempty"`
	Remote         string   `json:"remote,omitempty"`
	Experience     string   `json:"experience,omitempty"`
	Skills         string   `json:"skills,omitempty"`
	RequiredSkills []string `json:"required_skills,omitempty"`
	Paginatin      *PaginatedQuery
}
