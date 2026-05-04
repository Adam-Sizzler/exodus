package users

import (
	"testing"
	"time"
)

func stringPtr(value string) *string {
	return &value
}

func int64Ptr(value int64) *int64 {
	return &value
}

func TestPlannedUserStatusForUpdateReactivatesLimitedWhenLimitRaised(t *testing.T) {
	record := userRecord{
		Status:            "LIMITED",
		TrafficLimitBytes: 100,
	}
	req := updateUserRequest{TrafficLimitBytes: int64Ptr(200)}

	status, ok := plannedUserStatusForUpdate(record, req, time.Now())
	if !ok || status != "ACTIVE" {
		t.Fatalf("expected LIMITED user to reactivate when limit is raised, got ok=%v status=%q", ok, status)
	}
}

func TestPlannedUserStatusForUpdateReactivatesLimitedWhenLimitRemoved(t *testing.T) {
	record := userRecord{
		Status:            "LIMITED",
		TrafficLimitBytes: 100,
	}
	req := updateUserRequest{TrafficLimitBytes: int64Ptr(0)}

	status, ok := plannedUserStatusForUpdate(record, req, time.Now())
	if !ok || status != "ACTIVE" {
		t.Fatalf("expected LIMITED user to reactivate when limit is removed, got ok=%v status=%q", ok, status)
	}
}

func TestPlannedUserStatusForUpdateKeepsLimitedWhenLimitStillTooLow(t *testing.T) {
	record := userRecord{
		Status:            "LIMITED",
		TrafficLimitBytes: 100,
	}
	req := updateUserRequest{TrafficLimitBytes: int64Ptr(50)}

	status, ok := plannedUserStatusForUpdate(record, req, time.Now())
	if ok || status != "" {
		t.Fatalf("expected LIMITED user to stay limited, got ok=%v status=%q", ok, status)
	}
}

func TestPlannedUserStatusForUpdateReactivatesExpiredWhenExpireAtMovesToFuture(t *testing.T) {
	now := time.Date(2026, 5, 4, 12, 0, 0, 0, time.UTC)
	record := userRecord{
		Status:   "EXPIRED",
		ExpireAt: now.Add(-time.Hour),
	}
	req := updateUserRequest{ExpireAt: stringPtr(now.Add(24 * time.Hour).Format(time.RFC3339))}

	status, ok := plannedUserStatusForUpdate(record, req, now)
	if !ok || status != "ACTIVE" {
		t.Fatalf("expected EXPIRED user to reactivate when expiration moves to future, got ok=%v status=%q", ok, status)
	}
}

func TestPlannedUserStatusForUpdateExplicitStatusWins(t *testing.T) {
	now := time.Date(2026, 5, 4, 12, 0, 0, 0, time.UTC)
	record := userRecord{
		Status:   "EXPIRED",
		ExpireAt: now.Add(-time.Hour),
	}
	req := updateUserRequest{
		Status:   stringPtr("DISABLED"),
		ExpireAt: stringPtr(now.Add(24 * time.Hour).Format(time.RFC3339)),
	}

	status, ok := plannedUserStatusForUpdate(record, req, now)
	if !ok || status != "DISABLED" {
		t.Fatalf("expected explicit status to win, got ok=%v status=%q", ok, status)
	}
}

func TestUserConfigPresenceChanges(t *testing.T) {
	cases := []struct {
		name string
		prev string
		next string
		want bool
	}{
		{name: "active to disabled removes user", prev: "ACTIVE", next: "DISABLED", want: true},
		{name: "limited to active adds user", prev: "LIMITED", next: "ACTIVE", want: true},
		{name: "expired to limited stays absent", prev: "EXPIRED", next: "LIMITED", want: false},
		{name: "active to active unchanged", prev: "ACTIVE", next: "ACTIVE", want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := userConfigPresenceChanges(tc.prev, tc.next); got != tc.want {
				t.Fatalf("userConfigPresenceChanges(%q, %q) = %v, want %v", tc.prev, tc.next, got, tc.want)
			}
		})
	}
}
