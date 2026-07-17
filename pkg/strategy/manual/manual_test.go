package manual

import (
	"testing"
)

func TestValidateAndParse(t *testing.T) {
	tests := []struct {
		name    string
		version string
		strict  bool
		wantErr bool
		wantV   string
	}{
		// Valid versions - should pass in both modes
		{"valid semver", "1.2.3", false, false, "1.2.3"},
		{"valid semver prerelease", "1.2.3-alpha.1", false, false, "1.2.3-alpha.1"},
		{"valid semver metadata", "1.2.3+build.123", false, false, "1.2.3+build.123"},
		{"valid semver both", "1.2.3-beta+build.123", false, false, "1.2.3-beta+build.123"},
		{"valid semver strict", "1.2.3", true, false, "1.2.3"},
		{"valid semver prerelease strict", "1.2.3-alpha.1", true, false, "1.2.3-alpha.1"},

		// Partial versions - coerced in non-strict, rejected in strict
		{"partial major only non-strict", "1", false, false, "1.0.0"},
		{"partial major.minor non-strict", "1.2", false, false, "1.2.0"},
		{"partial major only strict", "1", true, true, ""},
		{"partial major.minor strict", "1.2", true, true, ""},

		// Leading zeros - accepted in non-strict, rejected in strict
		{"leading zero patch non-strict", "1.0.06", false, false, "1.0.6"},
		{"leading zero patch strict", "1.0.06", true, true, ""},
		{"date-style strict", "2026.06.29", true, true, ""},

		// Date-style versions - silently coerced in non-strict, rejected in strict
		{"date-style MM-DD-YYYY non-strict", "06-29-2026", false, false, "6.0.0-29-2026"},
		{"date-style YYYY-MM-DD non-strict", "2026-06-29", false, false, "2026.0.0-06-29"},
		{"date-style MM-DD-YYYY strict", "06-29-2026", true, true, ""},
		{"date-style YYYY-MM-DD strict", "2026-06-29", true, true, ""},

		// Invalid versions - rejected in both modes
		{"non-numeric leading segment", "stable-07-09-2026", false, true, ""},
		{"non-numeric leading segment strict", "stable-07-09-2026", true, true, ""},
		{"empty string", "", false, true, ""},
		{"empty string strict", "", true, true, ""},
		{"not a version", "hello", false, true, ""},
		{"not a version strict", "hello", true, true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := validateAndParse(tt.version, tt.strict)

			if (err != nil) != tt.wantErr {
				t.Errorf("validateAndParse(%q, strict=%v) error = %v, wantErr %v", tt.version, tt.strict, err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if v.String() != tt.wantV {
				t.Errorf("validateAndParse(%q, strict=%v) = %q, want %q", tt.version, tt.strict, v.String(), tt.wantV)
			}
		})
	}
}

func TestStrategy_ReadVersion(t *testing.T) {
	tests := []struct {
		name    string
		s       Strategy
		wantErr bool
		wantV   string
	}{
		{"valid non-strict", Strategy{Version: "1.2.3"}, false, "1.2.3"},
		{"valid strict", Strategy{Version: "1.2.3", Strict: true}, false, "1.2.3"},
		{"date-style non-strict", Strategy{Version: "06-29-2026"}, false, "6.0.0-29-2026"},
		{"date-style strict", Strategy{Version: "06-29-2026", Strict: true}, true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := tt.s.ReadVersion()

			if (err != nil) != tt.wantErr {
				t.Errorf("Strategy.ReadVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if v.String() != tt.wantV {
				t.Errorf("Strategy.ReadVersion() = %q, want %q", v.String(), tt.wantV)
			}
		})
	}
}
