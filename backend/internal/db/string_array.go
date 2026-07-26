package db

import (
	"encoding/json"
	"fmt"
	"strings"
)

// StringArray is a robust scanner for string arrays that supports JSON arrays
// (e.g. ["a","b"]) and PostgreSQL array literals (e.g. {"a","b"}).
type StringArray []string

// Scan implements sql.Scanner for StringArray.
func (a *StringArray) Scan(src any) error {
	if a == nil {
		return fmt.Errorf("StringArray: Scan on nil receiver")
	}

	switch v := src.(type) {
	case nil:
		*a = nil
		return nil
	case []string:
		out := make([]string, len(v))
		copy(out, v)
		*a = out
		return nil
	case string:
		return a.scanString(v)
	case []byte:
		return a.scanString(string(v))
	default:
		return fmt.Errorf("StringArray: unsupported source type %T", src)
	}
}

// Slice returns the underlying slice (or nil if empty).
func (a StringArray) Slice() []string {
	if a == nil {
		return nil
	}
	return []string(a)
}

func (a *StringArray) scanString(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		*a = []string{}
		return nil
	}
	if strings.HasPrefix(s, "[") {
		var out []string
		if err := json.Unmarshal([]byte(s), &out); err != nil {
			return fmt.Errorf("StringArray: parse json array: %w", err)
		}
		*a = out
		return nil
	}
	if strings.HasPrefix(s, "{") {
		out, err := parsePostgresTextArray(s)
		if err != nil {
			return fmt.Errorf("StringArray: parse postgres array: %w", err)
		}
		*a = out
		return nil
	}

	// Fallback: treat as single string value.
	*a = []string{s}
	return nil
}

// parsePostgresTextArray parses a PostgreSQL text array literal like {"a","b"}.
func parsePostgresTextArray(s string) ([]string, error) {
	s = strings.TrimSpace(s)
	if s == "{}" {
		return []string{}, nil
	}
	if len(s) < 2 || s[0] != '{' || s[len(s)-1] != '}' {
		return nil, fmt.Errorf("invalid array literal")
	}

	body := s[1 : len(s)-1]
	if body == "" {
		return []string{}, nil
	}

	var out []string
	var buf strings.Builder
	inQuotes := false
	escaped := false
	elementQuoted := false

	flush := func() {
		val := buf.String()
		if !elementQuoted {
			val = strings.TrimSpace(val)
			if val == "NULL" {
				val = ""
			}
		}
		out = append(out, val)
		buf.Reset()
		elementQuoted = false
	}

	for i := 0; i < len(body); i++ {
		ch := body[i]

		if inQuotes {
			if escaped {
				buf.WriteByte(ch)
				escaped = false
				continue
			}
			switch ch {
			case '\\':
				escaped = true
			case '"':
				inQuotes = false
			default:
				buf.WriteByte(ch)
			}
			continue
		}

		switch ch {
		case '"':
			if buf.Len() == 0 {
				inQuotes = true
				elementQuoted = true
				continue
			}
			buf.WriteByte(ch)
		case ',':
			flush()
		default:
			buf.WriteByte(ch)
		}
	}

	if inQuotes {
		return nil, fmt.Errorf("unterminated quoted string in array literal")
	}
	flush()
	return out, nil
}
