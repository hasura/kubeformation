package v1

// ClusterSpec defines the structure of a Kubernetes cluster
type ClusterSpec struct {
	// Version is the spec version
	Version string `json:"version"`

	// Name of the Kubernetes cluster
	Name string `json:"name"`

	// Provider for which the template has to be generated
	Provider string `json:"provider"`

	// K8SVersion is the version for Kubernetes Cluster. Specific version numbers
	// will vary from provider to provider
	K8SVersion string `json:"k8sVersion"`

	// NodePools denotes the node pools for the cluster
	NodePools []NodePool `json:"nodePools"`

	Volumes []Volume `json:"volumes"`
}

// NodePool indicated the spec for a Kubernetes node pools.
// A NodePool is just a collection of nodes having the same configuration
type NodePool struct {
	// Name of the node pool
	Name string `json:"name"`

	// Size is the number of nodes in this node pool
	Size int64 `json:"size"`

	// Type of the nodes in this pool. Value will be provider specific
	Type string `json:"type"`

	// OSImage is the OS to be used for the node
	OSImage string `json:"osImage"`

	// Labels to be applied to the nodes in the pool
	Labels map[string]string `json:"labels"`
}

type Volume struct {
	Name string `json:"name"`
	Size int    `json:"size"`
}
