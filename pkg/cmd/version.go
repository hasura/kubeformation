package cmd

var (
	// Version of the CLI/API, set during build time
	version = "v0.0.0-unset"
)

// GetVersion returns the current version string
func GetVersion() string {
	return version
}
