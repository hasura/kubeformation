package spec

// VersionedSpec denotes the minimum required information to identify a spec
// version. Cluster Spec in read as this type first to infer the version and
// choose the correct handler.
type VersionedSpec struct {
	// Version of the given spec
	Version string `json:"version"`
}
