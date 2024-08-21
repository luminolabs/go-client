package core

import "fmt"

// Version components
const (
	VersionMajor = 0       // Major version component of the current release
	VersionMinor = 1       // Minor version component of the current release
	VersionPatch = 0       // Patch version component of the current release
	VersionMeta  = "alpha" // Version metadata to append to the version string
)

// ClientName is the name of the Lumino client
const ClientName = "LuminoCLI"

// Version holds the textual version string.
var Version = func() string {
	return fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
}()

// VersionWithMeta holds the textual version string including the metadata.
var VersionWithMeta = func() string {
	v := Version
	if VersionMeta != "" {
		v += "-" + VersionMeta
	}
	return v
}()

// VersionWithName holds the textual version string including the client name.
var VersionWithName = func() string {
	return fmt.Sprintf("%s/%s", ClientName, VersionWithMeta)
}()

// GitCommit holds the Git commit hash, set during build time
var GitCommit string

// BuildDate holds the build date, set during build time
var BuildDate string

// VersionInfo returns a string with detailed version information
func VersionInfo() string {
	return fmt.Sprintf("%s\nVersion: %s\nGit Commit: %s\nBuild Date: %s",
		ClientName, VersionWithMeta, GitCommit, BuildDate)
}
