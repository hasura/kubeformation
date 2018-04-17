package spec

// VersionedSpec denotes the minimum required information to identify a spec version
type VersionedSpec struct {
	// Version of the given spec
	Version string `json:"version"`
}
