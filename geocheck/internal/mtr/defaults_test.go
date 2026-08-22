package mtr

import "testing"

func ids(targets []Target) map[string]bool {
	out := make(map[string]bool, len(targets))
	for _, t := range targets {
		out[t.ID] = true
	}
	return out
}

// TestDefaultSetComposition pins which targets a plain `geocheck` traces.
// The default set is the one nobody chooses deliberately, so anything added to
// it is paid for by every user on every run; it is worth stating outright.
func TestDefaultSetComposition(t *testing.T) {
	have := ids(DefaultTargets())

	for _, id := range []string{
		"tiktok", "github", "fastly", "instagram", "telegram", "telegram_dc5",
		"anthropic", "steam", "quad9", "discord", "akamai",
	} {
		if !have[id] {
			t.Errorf("%q should be in the default set", id)
		}
	}

	// DC1 and DC3 sit in AS59930 with DC2's AS62041 neighbour DC4; tracing a
	// second address inside an AS already measured costs a probe and tells you
	// nothing new, so only one target per Telegram network is default.
	for _, id := range []string{"telegram_dc1", "telegram_dc3", "telegram_dc4"} {
		if have[id] {
			t.Errorf("%q should not be in the default set; it shares an AS with a target already there", id)
		}
	}
}

// TestTelegramTagKeepsAllFive makes sure narrowing the default set did not
// remove the ability to compare the data centres against each other.
func TestTelegramTagKeepsAllFive(t *testing.T) {
	got := Select("telegram")
	if len(got) != 5 {
		t.Fatalf("-T telegram selected %d targets, want all 5", len(got))
	}
	have := ids(got)
	for _, id := range []string{
		"telegram_dc1", "telegram", "telegram_dc3", "telegram_dc4", "telegram_dc5",
	} {
		if !have[id] {
			t.Errorf("-T telegram is missing %q", id)
		}
	}
}

// TestGitHubTargetMatchesWhereItIsServedFrom records a measured fact that looks
// like a mistake: github.com resolves into Microsoft's AS8075, not GitHub's own
// AS36459, because the site is served from Azure. The expected ASN is what the
// path is judged against, so restoring 36459 would make every run report the
// destination as the wrong network.
func TestGitHubTargetMatchesWhereItIsServedFrom(t *testing.T) {
	for _, target := range Catalog {
		if target.ID != "github" {
			continue
		}
		if target.ASN != 8075 {
			t.Errorf("github expects AS%d; it is served from Azure, AS8075", target.ASN)
		}
		if target.ICMPSilent {
			t.Error("the Azure address answers ICMP, so github should not be marked ICMPSilent")
		}
		return
	}
	t.Fatal("no github target in the catalogue")
}

// TestAnthropicTargetMatchesWhereItIsServedFrom records the same kind of
// measured fact as the GitHub case: claude.ai left Cloudflare's AS13335 for
// Anthropic's own AS399358. Judging the path against 13335 would report every
// run as ending in the wrong network.
func TestAnthropicTargetMatchesWhereItIsServedFrom(t *testing.T) {
	for _, target := range Catalog {
		if target.ID != "anthropic" {
			continue
		}
		if target.ASN != 399358 {
			t.Errorf("anthropic expects AS%d; claude.ai is served from AS399358", target.ASN)
		}
		return
	}
	t.Fatal("no anthropic target in the catalogue")
}
