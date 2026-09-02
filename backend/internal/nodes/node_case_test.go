package users

import (
	"strings"
	"testing"
)

func TestCaseInsensitiveUserLookupResolution(t *testing.T) {
	// Simulate: User in DB is "Pusik" (or "PUsIk")
	dbUsername := "Pusik"
	var dbUserID int64 = 42

	// Node reports traffic under different casing: "pusik", "PUSIK", "PUsIk"
	incomingVariants := []string{"pusik", "PUSIK", "PUsIk", "Pusik"}

	for _, incoming := range incomingVariants {
		userIDs := make(map[string]int64)
		userIDsLower := make(map[string]int64)

		// 1. If exact match:
		if incoming == dbUsername {
			userIDs[dbUsername] = dbUserID
			userIDsLower[strings.ToLower(dbUsername)] = dbUserID
		} else {
			// Fallback case-insensitive match:
			if strings.ToLower(incoming) == strings.ToLower(dbUsername) {
				userIDs[dbUsername] = dbUserID
				userIDsLower[strings.ToLower(dbUsername)] = dbUserID
			}
		}

		// Resolution logic from stream.go:
		userID, ok := userIDs[incoming]
		if !ok {
			userID, ok = userIDsLower[strings.ToLower(incoming)]
		}

		if !ok || userID != dbUserID {
			t.Fatalf("Failed to resolve user ID for incoming variant %q: got ok=%v, userID=%d, want %d", incoming, ok, userID, dbUserID)
		}
	}
}
