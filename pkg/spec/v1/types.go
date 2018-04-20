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

	// Volumes are the persistent volumes to be used in this cluster.
	// Corresponding disk and k8s persistent volume object will be created.
	Volumes []Volume `json:"volumes"`
}

// NodePools indicate the spec for a Kubernetes node pools.
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

// Volume denotes a volume to be used in the cluster. A backing disk and a k8s
// PV object will be created for each volume.
type Volume struct {
	// Name of the volume
	Name string `json:"name"`

	// Size of the volume in GB
	Size int `json:"size"`
}
