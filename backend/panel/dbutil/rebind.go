package dbutil

import "strconv"

// Rebind converts SQL placeholders from '?' to PostgreSQL-style '$1, $2, ...'.
// It skips replacements inside single or double quoted strings.
func Rebind(query string) string {
	if query == "" {
		return query
	}

	out := make([]byte, 0, len(query)+8)
	arg := 1
	inSingle := false
	inDouble := false

	for i := 0; i < len(query); i++ {
		ch := query[i]

		switch ch {
		case '\'':
			if !inDouble {
				inSingle = !inSingle
			}
		case '"':
			if !inSingle {
				inDouble = !inDouble
			}
		case '?':
			if !inSingle && !inDouble {
				out = append(out, '$')
				out = strconv.AppendInt(out, int64(arg), 10)
				arg++
				continue
			}
		}

		out = append(out, ch)
	}

	return string(out)
}
