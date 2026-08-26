package util

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

// IsUniqueViolation reports whether err is a Postgres unique-constraint
// violation (SQLSTATE 23505). It checks the driver's typed *pgconn.PgError
// first — this is the single source of truth that replaces four different,
// hand-copied checks that used to live in internal/httpapi/users,
// subscriptiontemplate, subscriptionpageconfigs, and externalsquads, one of
// which (externalsquads) never checked pgconn.PgError at all and relied
// purely on substring-matching err.Error(), which is fragile: it depends on
// the exact wording the driver happens to produce.
//
// If constraintNames are given, the violation must match one of them
// (checked against pgErr.ConstraintName); with no names given, any unique
// violation counts. A substring-match fallback against err.Error() is kept
// for non-pgconn error paths (e.g. a test double that returns a plain error
// with a SQLite-style "UNIQUE constraint failed" message) — this mirrors
// what the strongest of the four original implementations already did.
func IsUniqueViolation(err error, constraintNames ...string) bool {
	if err == nil {
		return false
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		if len(constraintNames) == 0 {
			return true
		}
		for _, name := range constraintNames {
			if pgErr.ConstraintName == name {
				return true
			}
		}
		return false
	}

	message := err.Error()
	if len(constraintNames) == 0 {
		return strings.Contains(message, "duplicate key") ||
			strings.Contains(message, "Unique constraint") ||
			strings.Contains(message, "UNIQUE constraint failed")
	}
	for _, name := range constraintNames {
		if strings.Contains(message, name) {
			return true
		}
	}
	return false
}
