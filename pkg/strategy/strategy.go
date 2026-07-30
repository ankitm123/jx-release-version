package strategy

import (
	"strings"

	"github.com/Masterminds/semver/v3"
)

type VersionReader interface {
	ReadVersion() (*semver.Version, error)
}

type VersionBumper interface {
	BumpVersion(previous semver.Version, tagSuffix string) (*semver.Version, error)
}

func IncPatch(previous semver.Version, tagSuffix string) semver.Version {
	if tagSuffix != "" && previous.Prerelease() == strings.TrimPrefix(tagSuffix, "-") {
		// IncPatch doesn't increase version if pre release is set, instead it is cleared.
		// Since the latest version has the same pre release as would be set here we would end up
		// trying to change the previous tag.
		// Instead the first IncPatch clears pre release before the actual increment.
		previous = previous.IncPatch()
	}
	return previous.IncPatch()
}
