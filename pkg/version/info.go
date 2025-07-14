package version

import (
	"fmt"
	"runtime/debug"
)

var (
	Version = "unknown"
	Commit  = "unknown"
	Date    = "unknown"
)

func String() string {
	versionString := fmt.Sprintf("%s (built %s)", Version, Date)
	if buildInfo, available := debug.ReadBuildInfo(); available {
		versionString = fmt.Sprintf("%s (built %s with %s)", Version, Date, buildInfo.GoVersion)
	}
	return versionString
}
