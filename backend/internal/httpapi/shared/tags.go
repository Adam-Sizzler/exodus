package shared

import "strings"

type SetEntityTagsRequest struct {
	UUID string   `json:"uuid"`
	Tags []string `json:"tags"`
}

type EntityTagsResponse struct {
	Tags []string `json:"tags"`
}

type SetEntityTagsResponse struct {
	UUID string   `json:"uuid"`
	Tags []string `json:"tags"`
}

func SanitizeTags(raw []string) []string {
	seen := make(map[string]struct{}, len(raw))
	res := make([]string, 0, len(raw))
	for _, t := range raw {
		trimmed := strings.TrimSpace(t)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; !ok {
			seen[trimmed] = struct{}{}
			res = append(res, trimmed)
		}
	}
	return res
}
