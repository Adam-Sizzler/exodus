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
			want: []string{"DAY"},
		},
		{
			name: "weekly reset window",
			now:  time.Date(2026, 5, 4, 0, 8, 0, 0, time.Local),
			want: []string{"WEEK"},
		},
		{
			name: "monthly reset window",
			now:  time.Date(2026, 6, 1, 0, 10, 0, 0, time.Local),
			want: []string{"MONTH"},
		},
		{
			name: "outside reset windows",
			now:  time.Date(2026, 5, 4, 1, 0, 0, 0, time.Local),
			want: []string{},
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
