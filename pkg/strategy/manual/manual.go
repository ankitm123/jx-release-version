package manual

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
)

type Strategy struct {
	Version string
	Strict  bool
}

func (s Strategy) ReadVersion() (*semver.Version, error) {
	log.Logger().Debugf("Using manual version %s (strict: %v)", s.Version, s.Strict)
	return validateAndParse(s.Version, s.Strict)
}

func (s Strategy) BumpVersion(_ semver.Version) (*semver.Version, error) {
	log.Logger().Debugf("Using manual version %s (strict: %v)", s.Version, s.Strict)
	return validateAndParse(s.Version, s.Strict)
}

// validateAndParse parses a version string using strict or relaxed semver validation.
func validateAndParse(version string, strict bool) (*semver.Version, error) {
	var v *semver.Version
	var err error

	if strict {
		v, err = semver.StrictNewVersion(version)
	} else {
		v, err = semver.NewVersion(version)
	}

	if err != nil {
		return nil, fmt.Errorf("invalid semantic version %q: %w", version, err)
	}

	return v, nil
}
