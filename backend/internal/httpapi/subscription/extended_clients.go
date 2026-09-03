package subscription

import (
	"regexp"
	"sync"
)

var extendedClientsRegexes = []*regexp.Regexp{
	regexp.MustCompile(`^FlClash ?X/`),
	regexp.MustCompile(`^Flowvy/`),
	regexp.MustCompile(`^prizrak-box/`),
	regexp.MustCompile(`^koala-clash/`),
	regexp.MustCompile(`^Happ/`),
	regexp.MustCompile(`^INCY/`),
}

var jsonSubscriptionFallbackRegexes = []*regexp.Regexp{
	regexp.MustCompile(`^[Ss]treisand`),
	regexp.MustCompile(`^Happ/`),
	regexp.MustCompile(`^INCY/`),
	regexp.MustCompile(`^ktor-client`),
	regexp.MustCompile(`^V2Box`),
	regexp.MustCompile(`^io\.github\.saeeddev94\.xray/`),
	regexp.MustCompile(`^v2rayNG/(\d+\.\d+\.\d+)`),
	regexp.MustCompile(`^v2rayN/(\d+\.\d+\.\d+)`),
	regexp.MustCompile(`^v2plus/(\d+\.\d+\.\d+)`),
}

var (
	ruleRegexCacheMu sync.RWMutex
	ruleRegexCache   = make(map[string]*regexp.Regexp)
)

// isExtendedClient checks if the client user-agent belongs to extended clients
// (Happ, INCY, FlClashX, Flowvy, koala-clash, prizrak-box) either by built-in regexes
// or by custom regexes specified in the matched SRR rule.
func isExtendedClient(userAgent string, additionalPatterns []string) bool {
	if userAgent == "" {
		return false
	}
	for _, re := range extendedClientsRegexes {
		if re.MatchString(userAgent) {
			return true
		}
	}
	for _, pattern := range additionalPatterns {
		if pattern == "" {
			continue
		}
		re := getCachedRegex(pattern)
		if re != nil && re.MatchString(userAgent) {
			return true
		}
	}
	return false
}

// isJSONSubscriptionFallbackSupported checks if the client user-agent can accept
// a JSON subscription format (XRAY_JSON) instead of raw base64 links when
// serveJsonAtBaseSubscription is enabled.
func isJSONSubscriptionFallbackSupported(userAgent string) bool {
	if userAgent == "" {
		return false
	}
	for _, re := range jsonSubscriptionFallbackRegexes {
		if re.MatchString(userAgent) {
			return true
		}
	}
	return false
}

func getCachedRegex(pattern string) *regexp.Regexp {
	ruleRegexCacheMu.RLock()
	re, ok := ruleRegexCache[pattern]
	ruleRegexCacheMu.RUnlock()
	if ok {
		return re
	}

	ruleRegexCacheMu.Lock()
	defer ruleRegexCacheMu.Unlock()
	if re, ok := ruleRegexCache[pattern]; ok {
		return re
	}

	compiled, err := regexp.Compile(pattern)
	if err != nil {
		ruleRegexCache[pattern] = nil
		return nil
	}
	ruleRegexCache[pattern] = compiled
	return compiled
}
