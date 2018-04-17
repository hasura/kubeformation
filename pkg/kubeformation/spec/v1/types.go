package v1

// ClusterSpec defines the structure of a Kubernetes cluster
type ClusterSpec struct {
	// Version is the spec version
	Version string `json:"version"`

	// Name of the Kubernetes cluster
	Name string `json:"name"`

	// Provider for which the template has to be generated
	Provider string `json:"provider"`

	// K8SVersion is the version for Kubernetes Cluster
	K8SVersion string `json:"k8sVersion"`

	// Nodes denotes the node pools for the cluster
	Nodes []NodePool `json:"nodes"`
}

// NodePool indicated the spec for a Kubernetes node pools.
// A NodePool is just a collection of nodes having the same configuration
type NodePool struct {
	// Name of the node pool
	Name string `json:"name"`

	// Labels to be applied to the nodes in the pool
	Labels map[string]string `json:"labels"`

	// PoolSize is the number of nodes in this node pool
	PoolSize int64 `json:"poolSize"`
}
