// Package jsonx provides just enough JSON path extraction for scraping the
// small, well-known response shapes the geo providers return.
package jsonx

import (
	"encoding/json"
	"strconv"
	"strings"
)

// Get walks a dot-separated path through decoded JSON and returns the value at
// the end. Numeric segments index into arrays. A missing path yields nil.
//
//	Get(body, "location.country.code")
//	Get(body, "0.data.requestInfo.countryCode")
func Get(data []byte, path string) any {
	var root any
	if err := json.Unmarshal(data, &root); err != nil {
		return nil
	}
	return walk(root, path)
}

func walk(node any, path string) any {
	if path == "" {
		return node
	}
	for _, seg := range strings.Split(path, ".") {
		if node == nil {
			return nil
		}
		switch v := node.(type) {
		case map[string]any:
			node = v[seg]
		case []any:
			i, err := strconv.Atoi(seg)
			if err != nil || i < 0 || i >= len(v) {
				return nil
			}
			node = v[i]
		default:
			return nil
		}
	}
	return node
}

// String extracts a path as a string. Numbers and booleans are stringified;
// JSON null and missing paths yield "".
func String(data []byte, path string) string {
	switch v := Get(data, path).(type) {
	case nil:
		return ""
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return ""
	}
}

// Valid reports whether data parses as JSON.
func Valid(data []byte) bool { return json.Valid(data) }
