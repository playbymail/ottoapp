// Package ottoapp implements the OttoMap web server.
package ottoapp

import (
	"github.com/maloquacious/semver"
)

var (
	version = semver.Version{Major: 0, Minor: 10, Patch: 0, PreRelease: "alpha", Build: semver.Commit()}
)

func Version() semver.Version {
	return version
}
