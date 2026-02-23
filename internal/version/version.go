// Package version provides build-time version information.
//
//nolint:revive // "version" is not actually a standard library package
package version

// Get returns the version string.
func Get() string {
	return getVersion()
}

// GetCommit returns the commit hash.
func GetCommit() string {
	return getCommit()
}

// GetDate returns the build date.
func GetDate() string {
	return getDate()
}

// Info returns the build version information.
func Info() (string, string, string) {
	return getVersion(), getCommit(), getDate()
}
