// Package version carries the build identity, stamped at link time.
package version

import (
	"runtime/debug"

	"exodus-node/constant"
)

var (
	// Version is the release tag, matching node version.
	Version = constant.Version
	// Commit is the git revision.
	Commit = ""
	// Date is the build timestamp.
	Date = ""
)

// String renders a short human-readable version.
func String() string {
	v := Version
	if v == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, s := range info.Settings {
				if s.Key == "vcs.revision" && len(s.Value) >= 7 {
					return v + " (" + s.Value[:7] + ")"
				}
			}
		}
		return v
	}
	if Commit != "" && len(Commit) >= 7 {
		return v + " (" + Commit[:7] + ")"
	}
	return v
}
