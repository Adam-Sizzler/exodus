package userwatchdog

import (
	"testing"
	"time"
)

func TestScheduledResetStrategies(t *testing.T) {
	cases := []struct {
		name string
		now  time.Time
		want []string
	}{
		{
			name: "daily reset window",
			now:  time.Date(2026, 5, 5, 0, 3, 0, 0, time.Local),
			want: []string{"DAY", "WEEK", "MONTH"},
		},
		{
			name: "weekly reset window",
			now:  time.Date(2026, 5, 4, 0, 8, 0, 0, time.Local),
			want: []string{"DAY", "WEEK", "MONTH"},
		},
		{
			name: "monthly reset window",
			now:  time.Date(2026, 6, 1, 0, 10, 0, 0, time.Local),
			want: []string{"DAY", "WEEK", "MONTH"},
		},
		{
			name: "outside old reset windows",
			now:  time.Date(2026, 5, 4, 1, 0, 0, 0, time.Local),
			want: []string{"DAY", "WEEK", "MONTH"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := scheduledResetStrategies(tc.now)
			if len(got) != len(tc.want) {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("got %v, want %v", got, tc.want)
				}
			}
		})
	}
}

func TestResetPeriodBoundary(t *testing.T) {
	loc := time.FixedZone("test", 7*60*60)
	now := time.Date(2026, 5, 7, 13, 45, 30, 0, loc)

	cases := []struct {
		strategy string
		want     time.Time
		ok       bool
	}{
		{
			strategy: "DAY",
			want:     time.Date(2026, 5, 7, 0, 0, 0, 0, loc),
			ok:       true,
		},
		{
			strategy: "WEEK",
			want:     time.Date(2026, 5, 4, 0, 0, 0, 0, loc),
			ok:       true,
		},
		{
			strategy: "MONTH",
			want:     time.Date(2026, 5, 1, 0, 0, 0, 0, loc),
			ok:       true,
		},
		{
			strategy: "NO_RESET",
			want:     time.Time{},
			ok:       false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.strategy, func(t *testing.T) {
			got, ok := resetPeriodBoundary(tc.strategy, now)
			if ok != tc.ok {
				t.Fatalf("ok=%t, want %t", ok, tc.ok)
			}
			if !got.Equal(tc.want) {
				t.Fatalf("got %s, want %s", got, tc.want)
			}
		})
	}
}
