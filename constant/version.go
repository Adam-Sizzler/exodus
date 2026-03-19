package constant

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

var (
	Version    = "unknown"
	Revision   = "unknown"
	BuildTags  = "unknown"
	CgoEnabled = "unknown"
)

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
		"cerberus version: %s\n"+
			"environment: %s %s\n"+
			"tags: %s\n"+
			"revision: %s%s\n"+
			"cgo: %s",
		Version, goVersion, goEnv, BuildTags, Revision, vcsDirty, CgoEnabled,
	)
}
