package users

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func parseUsersTableQuery(r *http.Request) (usersTableQuery, error) {
	values := r.URL.Query()
	query := usersTableQuery{
		Start:       0,
		Size:        25,
		Filters:     []usersTableFilter{},
		FilterModes: map[string]string{},
		Sorting:     []usersTableSorting{},
	}

	if raw := strings.TrimSpace(values.Get("start")); raw != "" {
		start, err := strconv.Atoi(raw)
		if err != nil {
			return query, err
		}
		if start > 0 {
			query.Start = start
		}
	}
	if raw := strings.TrimSpace(values.Get("size")); raw != "" {
		size, err := strconv.Atoi(raw)
		if err != nil {
			return query, err
		}
		if size > 0 {
			query.Size = size
		}
		if query.Size > 1000 {
			query.Size = 1000
		}
	}
	if raw := strings.TrimSpace(values.Get("filters")); raw != "" {
		if err := json.Unmarshal([]byte(raw), &query.Filters); err != nil {
			return query, err
		}
	}
	if raw := strings.TrimSpace(values.Get("filterModes")); raw != "" {
		if err := json.Unmarshal([]byte(raw), &query.FilterModes); err != nil {
			return query, err
		}
	}
	if raw := strings.TrimSpace(values.Get("sorting")); raw != "" {
		if err := json.Unmarshal([]byte(raw), &query.Sorting); err != nil {
			return query, err
		}
	}

	return query, nil
}

