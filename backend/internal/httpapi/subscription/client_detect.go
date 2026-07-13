package subscription

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func extractHwidHeaders(r *http.Request) *HwidHeaders {
	hwid := strings.TrimSpace(r.Header.Get("X-HWID"))
	if hwid == "" {
		return nil
	}
	userAgent := firstNonEmptyHeader(r, "User-Agent", "X-HWID-User-Agent")
	platform := firstNonEmptyLowerHeader(r, "X-Device-OS", "X-HWID-Platform")
	osVersion := firstNonEmptyHeader(r, "X-Ver-OS", "X-HWID-OS-Version")
	deviceModel := firstNonEmptyHeader(r, "X-Device-Model", "X-HWID-Device-Model")
	platform, osVersion, deviceModel, userAgent = normalizeHwidMetadata(platform, osVersion, deviceModel, userAgent)

	h := &HwidHeaders{
		Hwid:        hwid,
		Platform:    platform,
		OsVersion:   osVersion,
		DeviceModel: deviceModel,
		UserAgent:   userAgent,
	}
	return h
}

func extractSyntheticHwidHeaders(r *http.Request, userUUID, requestIP string) *HwidHeaders {
	userAgent := strings.TrimSpace(r.Header.Get("User-Agent"))
	platform := firstNonEmptyLowerHeader(r, "X-Device-OS", "X-HWID-Platform")
	osVersion := firstNonEmptyHeader(r, "X-Ver-OS", "X-HWID-OS-Version")
	deviceModel := firstNonEmptyHeader(r, "X-Device-Model", "X-HWID-Device-Model")

	hasMetadata := userAgent != "" || platform != nil || osVersion != nil || deviceModel != nil
	if !hasMetadata {
		return nil
	}
	platform, osVersion, deviceModel, userAgentPtr := normalizeHwidMetadata(
		platform,
		osVersion,
		deviceModel,
		stringPtrIfNotEmpty(userAgent),
	)

	signature := strings.Join([]string{
		"exodus:synthetic-hwid:v1",
		"ua=" + strings.ToLower(ptrString(userAgentPtr)),
		"platform=" + strings.ToLower(ptrString(platform)),
		"os=" + strings.ToLower(ptrString(osVersion)),
		"model=" + strings.ToLower(ptrString(deviceModel)),
	}, "|")

	return &HwidHeaders{
		Hwid:        deterministicSyntheticHwid(userUUID, signature),
		Platform:    platform,
		OsVersion:   osVersion,
		DeviceModel: deviceModel,
		UserAgent:   userAgentPtr,
		RequestIP:   stringPtrIfNotEmpty(requestIP),
		Synthetic:   true,
	}
}

func normalizeHwidMetadata(platform, osVersion, deviceModel, userAgent *string) (*string, *string, *string, *string) {
	normalizedUserAgent := stringPtrIfNotEmpty(ptrString(userAgent))
	normalizedPlatform := lowerStringPtr(platform)
	if normalizedPlatform == nil {
		if inferred := inferPlatformFromUserAgent(ptrString(normalizedUserAgent)); inferred != "" {
			normalizedPlatform = &inferred
		}
	}
	if deviceModel == nil {
		deviceModel = stringPtrIfNotEmpty("unknown")
	}
	return normalizedPlatform, stringPtrIfNotEmpty(ptrString(osVersion)), stringPtrIfNotEmpty(ptrString(deviceModel)), normalizedUserAgent
}

func deterministicSyntheticHwid(userUUID, signature string) string {
	namespace, err := uuid.Parse(strings.TrimSpace(userUUID))
	if err != nil {
		namespace = uuid.NameSpaceOID
	}
	return uuid.NewSHA1(namespace, []byte(signature)).String()
}

func firstNonEmptyHeader(r *http.Request, names ...string) *string {
	for _, name := range names {
		value := strings.TrimSpace(r.Header.Get(name))
		if value == "" {
			continue
		}
		return &value
	}
	return nil
}

func firstNonEmptyLowerHeader(r *http.Request, names ...string) *string {
	value := firstNonEmptyHeader(r, names...)
	if value == nil {
		return nil
	}
	return lowerStringPtr(value)
}

func stringPtrIfNotEmpty(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func lowerStringPtr(value *string) *string {
	if value == nil {
		return nil
	}
	lowered := strings.ToLower(strings.TrimSpace(*value))
	if lowered == "" {
		return nil
	}
	return &lowered
}

func ptrString(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func inferClientAppFromUserAgent(userAgent string) string {
	userAgent = strings.TrimSpace(userAgent)
	if userAgent == "" {
		return ""
	}
	for _, sep := range []string{"/", " ", "(", ";"} {
		if idx := strings.Index(userAgent, sep); idx > 0 {
			return strings.TrimSpace(userAgent[:idx])
		}
	}
	return userAgent
}

func inferKnownClientPlatform(client string) string {
	client = strings.ToLower(strings.TrimSpace(client))

	switch client {
	case "sfa",
		"sfatv",
		"sfandroidtv",
		"v2rayng",
		"exclave",
		"nekoboxforandroid",
		"matsuri",
		"sagernet",
		"clashforandroid",
		"clashmetaforandroid",
		"cmfa":
		return "android"

	case "sfi",
		"streisand",
		"v2box",
		"rabbithole",
		"shadowrocket":
		return "ios"

	case "sft":
		return "tvos"

	case "sfw",
		"v2rayn":
		return "windows"

	case "sfm",
		"v2rayu",
		"v2rayx",
		"v2rayxs",
		"clashx":
		return "macos"

	case "sfl":
		return "linux"

	default:
		return ""
	}
}

func inferPlatformFromUserAgent(userAgent string) string {
	lower := strings.ToLower(userAgent)
	if lower == "" {
		return ""
	}

	if platform := inferKnownClientPlatform(inferClientAppFromUserAgent(lower)); platform != "" {
		return platform
	}

	if idx := strings.Index(lower, "platform/"); idx >= 0 {
		rest := userAgent[idx+len("platform/"):]
		for i, r := range rest {
			if !(r == '-' || r == '_' || r == '.' || r == '/' || r >= '0' && r <= '9' || r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z') {
				rest = rest[:i]
				break
			}
		}
		if rest = strings.Trim(rest, "/ "); rest != "" {
			return strings.ToLower(rest)
		}
	}

	switch {
	case strings.Contains(lower, "windows"):
		return "windows"
	case strings.Contains(lower, "android"):
		return "android"
	case strings.Contains(lower, "iphone") || strings.Contains(lower, "ipad") || strings.Contains(lower, "ios"):
		return "ios"
	case strings.Contains(lower, "mac os") || strings.Contains(lower, "macos") || strings.Contains(lower, "macintosh") || strings.Contains(lower, "darwin"):
		return "macos"
	case strings.Contains(lower, "linux"):
		return "linux"
	default:
		return ""
	}
}
