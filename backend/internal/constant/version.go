package constant

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
)

var (
	Version        = "unknown"
	Revision       = "unknown"
	FrontendCommit = "unknown"
	GitBranch      = "unknown"
	BuildTime      = "unknown"
	BuildNumber    = "unknown"
	BuildTags      = "none"
	CgoEnabled     = "0"
)

func getEnvOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func GetVersion() string {
	return getEnvOrDefault("__EX_METADATA_VERSION", Version)
}

func GetBackendCommit() string {
	return getEnvOrDefault("__EX_METADATA_GIT_BACKEND_COMMIT", Revision)
}

func GetFrontendCommit() string {
	return getEnvOrDefault("__EX_METADATA_GIT_FRONTEND_COMMIT", FrontendCommit)
}

func GetGitBranch() string {
	return getEnvOrDefault("__EX_METADATA_GIT_BRANCH", GitBranch)
}

func GetBuildTime() string {
	return getEnvOrDefault("__EX_METADATA_BUILD_TIME", BuildTime)
}

func GetBuildNumber() string {
	return getEnvOrDefault("__EX_METADATA_BUILD_NUMBER", BuildNumber)
}

func GetBuildInfo() string {
	goVersion := runtime.Version()
	goEnv := runtime.GOOS + "/" + runtime.GOARCH

	vcsDirty := ""
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.modified" {
				if setting.Value == "true" {
					vcsDirty = " (dirty)"
				}
				break
			}
		}
	}

	return fmt.Sprintf(
		"exodus version: %s\n"+
			"branch: %s\n"+
			"build time: %s\n"+
			"build number: %s\n"+
			"backend commit: %s%s\n"+
			"frontend commit: %s\n"+
			"environment: %s %s\n"+
			"tags: %s\n"+
			"cgo: %s",
		GetVersion(), GetGitBranch(), GetBuildTime(), GetBuildNumber(), GetBackendCommit(), vcsDirty, GetFrontendCommit(), goVersion, goEnv, BuildTags, CgoEnabled,
	)
}
