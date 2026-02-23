// Package version provides build-time version information.
//
//nolint:revive // "version" is not actually a standard library package
package version

// These variables are set via ldflags during build.
var (
	versionValue = "dev"
	commitValue  = "unknown"
	dateValue    = "unknown"
)

// Get returns the version string.
func Get() string {
	return versionValue
}

// GetCommit returns the commit hash.
func GetCommit() string {
	return commitValue
}

// GetDate returns the build date.
func GetDate() string {
	return dateValue
}

// Info returns the build version information.
func Info() (string, string, string) {
	return versionValue, commitValue, dateValue
}
