package scheduler

import (
	"database/sql"
	"testing"
	"time"
)

func TestLatestNodeTrafficResetBoundary(t *testing.T) {
	location := time.FixedZone("UTC+7", 7*60*60)

	tests := []struct {
		name     string
		now      time.Time
		resetDay int
		want     time.Time
	}{
		{
			name:     "first day already passed",
			now:      time.Date(2026, time.June, 3, 1, 30, 0, 0, location),
			resetDay: 1,
			want:     time.Date(2026, time.June, 1, 1, 0, 0, 0, location),
		},
		{
			name:     "future day uses previous month",
			now:      time.Date(2026, time.June, 3, 1, 30, 0, 0, location),
			resetDay: 15,
			want:     time.Date(2026, time.May, 15, 1, 0, 0, 0, location),
		},
		{
			name:     "day beyond short month falls back to last day",
			now:      time.Date(2026, time.February, 28, 2, 0, 0, 0, location),
			resetDay: 31,
			want:     time.Date(2026, time.February, 28, 1, 0, 0, 0, location),
		},
		{
			name:     "before reset hour uses previous month",
			now:      time.Date(2026, time.February, 28, 0, 30, 0, 0, location),
			resetDay: 31,
			want:     time.Date(2026, time.January, 31, 1, 0, 0, 0, location),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := latestNodeTrafficResetBoundary(tt.now, tt.resetDay)
			if !got.Equal(tt.want) {
				t.Fatalf("boundary mismatch: got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestNodeTrafficResetDue(t *testing.T) {
	location := time.FixedZone("UTC+7", 7*60*60)
	now := time.Date(2026, time.June, 3, 1, 30, 0, 0, location)
	boundary := time.Date(2026, time.June, 1, 1, 0, 0, 0, location)

	tests := []struct {
		name        string
		createdAt   time.Time
		lastResetAt sql.NullTime
		wantDue     bool
	}{
		{
			name:      "due when no reset exists after boundary",
			createdAt: time.Date(2026, time.May, 1, 0, 0, 0, 0, location),
			wantDue:   true,
		},
		{
			name:      "not due when node was created after boundary",
			createdAt: time.Date(2026, time.June, 2, 0, 0, 0, 0, location),
			wantDue:   false,
		},
		{
			name:      "not due when reset already happened at boundary",
			createdAt: time.Date(2026, time.May, 1, 0, 0, 0, 0, location),
			lastResetAt: sql.NullTime{
				Time:  boundary,
				Valid: true,
			},
			wantDue: false,
		},
		{
			name:      "not due when reset already happened after boundary",
			createdAt: time.Date(2026, time.May, 1, 0, 0, 0, 0, location),
			lastResetAt: sql.NullTime{
				Time:  boundary.Add(time.Hour),
				Valid: true,
			},
			wantDue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBoundary, gotDue := nodeTrafficResetDue(now, 1, tt.createdAt, tt.lastResetAt)
			if !gotBoundary.Equal(boundary) {
				t.Fatalf("boundary mismatch: got %s, want %s", gotBoundary, boundary)
			}
			if gotDue != tt.wantDue {
				t.Fatalf("due mismatch: got %t, want %t", gotDue, tt.wantDue)
			}
		})
	}
}
