// Package ottoapp implements the OttoMap web server.
package ottoapp

import (
	"github.com/maloquacious/semver"
)

var (
	version = semver.Version{
		Major: 0,
		Minor: 44,
		Patch: 3,
		Build: semver.Commit(),
	}
)

func Version() semver.Version {
	return version
}
