package shared

import (
	"encoding/json"
	"fmt"
	"strings"
)

// PostgresTextArrayLiteral formats a Go string slice into a PostgreSQL text[] array literal,
// e.g. {"a","b","c"} with proper escaping of quotes and backslashes.
func PostgresTextArrayLiteral(items []string) string {
	if len(items) == 0 {
		return "{}"
	}
	quoted := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		item = strings.ReplaceAll(item, `\`, `\\`)
		item = strings.ReplaceAll(item, `"`, `\"`)
		quoted = append(quoted, fmt.Sprintf(`"%s"`, item))
	}
	return "{" + strings.Join(quoted, ",") + "}"
}

// ParsePgTextArray parses either a PostgreSQL text[] array literal (e.g. "{a,b,c}")
// or a JSON array string (e.g. `["a","b","c"]`) into a Go string slice.
func ParsePgTextArray(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "{}" || raw == "[]" {
		return []string{}
	}
	if strings.HasPrefix(raw, "[") && strings.HasSuffix(raw, "]") {
		var list []string
		if err := json.Unmarshal([]byte(raw), &list); err == nil && list != nil {
			return list
		}
	}
	raw = strings.TrimPrefix(raw, "{")
	raw = strings.TrimSuffix(raw, "}")
	if raw == "" {
		return []string{}
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(strings.Trim(p, `"`))
		if p == "" {
			continue
		}
		out = append(out, p)
	}
	return out
}
