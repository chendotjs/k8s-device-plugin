package version

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	// commitHash contains the current Git revision. Use make to build to make
	// sure this gets set.
	CommitHash string

	// buildDate contains the date of the current build.
	BuildDate string
)

// Version represents the aios-operator build version.
type Version struct {
	// Major and minor version.
	Number float32

	// Increment this for bug releases
	PatchLevel uint

	// aios-operator VersionSuffix is the suffix used in the aios-operator version string.
	// It will be blank for release versions.
	Suffix string
}

func (v Version) String() string {
	return version(v.Number, v.PatchLevel, v.Suffix)
}

// Version returns the aios-operator version.
func (v Version) Version() VersionString {
	return VersionString(v.String())
}

// VersionString represents a aios-operator version string.
type VersionString string

func (h VersionString) String() string {
	return string(h)
}

// BuildVersionString creates a version string. This is what you see when
// running "aios-operator version".
func BuildVersionString() string {
	program := "k8s-device-plugin"

	version := "v" + CurrentVersion.String()
	if CommitHash != "" {
		version += "-" + strings.ToUpper(CommitHash)
	}

	osArch := runtime.GOOS + "/" + runtime.GOARCH

	date := BuildDate
	if date == "" {
		date = "unknown"
	}

	return fmt.Sprintf("%s %s %s BuildDate: %s", program, version, osArch, date)

}

func version(version float32, patchVersion uint, suffix string) string {
	return fmt.Sprintf("%.g.%d%s", version, patchVersion, suffix)
}
