package detect

import "testing"

// The answers below are real, captured from three hosts: two whose DNS reaches
// OpenDNS and one whose port 53 is intercepted. The identifiers name OpenDNS
// sites, not anybody's location.
func TestOpenDNSFinding(t *testing.T) {
	cases := []struct {
		name    string
		answers []string
		want    bool
	}{
		{
			// What a genuine OpenDNS resolver returns: its own server id. This
			// is the case the previous implementation reported as interception,
			// which made the check fire on every healthy host.
			name:    "a server identifier means the query arrived",
			answers: []string{"r3001.lon"}, want: false,
		},
		{
			name:    "another site, same shape",
			answers: []string{"r3007.ams"}, want: false,
		},
		{
			// The authoritative record. Reaching it means the query was
			// answered by a resolver that is not OpenDNS.
			name:    "the authoritative sentence means it did not",
			answers: []string{"I am not an OpenDNS resolver."}, want: true,
		},
		{
			name:    "case and spacing do not matter",
			answers: []string{"i am NOT an opendns resolver."}, want: true,
		},
		{
			name:    "no answers proves nothing",
			answers: nil, want: false,
		},
		{
			// Interception can leave several records; one giveaway is enough.
			name:    "any answer carrying the sentence is enough",
			answers: []string{"r3001.lon", "I am not an OpenDNS resolver."}, want: true,
		},
		{
			// Unrecognised, but not evidence of interception either.
			name:    "something else is not asserted on",
			answers: []string{"unexpected"}, want: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := openDNSFinding(c.answers)
			if (got != nil) != c.want {
				t.Fatalf("finding=%v, want %v", got != nil, c.want)
			}
			if got != nil && got.Severity != Alert {
				t.Errorf("severity = %v, want Alert", got.Severity)
			}
		})
	}
}

// TestOpenDNSCheckCanPass is the regression that matters. The check it replaces
// looked for a sentence no resolver returns, so it could only ever fire. A
// check that cannot come back clean is not a check.
func TestOpenDNSCheckCanPass(t *testing.T) {
	if openDNSFinding([]string{"r3001.lon"}) != nil {
		t.Fatal("a healthy answer still produced a finding")
	}
}
