// Package gke implements GKE provider for Kubeformation. It can generate Google
// Deployment Manager template in Jinja format. Only input required is the Spec
// object.
//
// Google Deployment Manager template supports Python files, Jinja templates as
// well as plain Yaml files. Kubeformation converts cluster spec into "Spec" as
// defined below and generate two files, cluster.jinja and cluster.yaml. These
// files are defined as Go templates in templates.go and are rendered with Spec
// as context.
package gke

import (
	"bytes"
	"html/template"

	"github.com/hasura/kubeformation/pkg/provider"
)

const (
	// Default version to be used in case cluster spec did not contain k8sVersion.
	DefaultK8SVersion = "1.9"

	// Google compute machine type to be used as default.
	DefaultMachineType = "n1-standard-1"

	// Default value for OS image.
	DefaultImageType = "cos"
)

// funcMap is a template helper function
var funcMap = template.FuncMap{
	"sub": func(i int) int {
		if i == 0 {
			return 0
		}
		return i - 1
	},
}

// Spec defines the context required to render GDM template.
type Spec struct {
	// Name of the cluster
	Name string

	// Kubernetes version to use
	K8SVersion string

	// Node pools for the cluster
	NodePools []NodePool

	// Persistent disks to be created
	Volumes []Volume
}

// NodePool defines a collection of nodes with similar properties in a
// Kubernetes cluster.
type NodePool struct {
	// Name of the node pool
	Name string

	// Kubernetes labels to be applied to nodes in this pool
	Labels map[string]string

	// Number of nodes in this pool
	Size int64

	// Google compute engine machine type to use
	MachineType string

	// Image type to use for the nodes
	ImageType string
}

// Volume defines a Google Persistent Disk to be created along with the cluster.
type Volume struct {
	// Name of the disk
	Name string

	// Size of the disk to be created, in GB
	SizeGB int
}

// NewDefaultSpec returns a spec object which is complete enough to render a
// valid GDM template.
func NewDefaultSpec() *Spec {
	return &Spec{
		Name:       "gke-cluster",
		K8SVersion: DefaultK8SVersion,
		NodePools: []NodePool{
			NodePool{
				Name:        "np-1",
				Size:        1,
				MachineType: DefaultMachineType,
				ImageType:   DefaultImageType,
			},
		},
	}
}

// GetType returns the type of this provider.
func (s *Spec) GetType() provider.ProviderType {
	return provider.GKE
}

// MarshalFiles returns the rendered GDM template as a map which indicates
// filename as keys and the file content as value.
// FIXME: test does not capture the template errors.
func (s *Spec) MarshalFiles() (map[string][]byte, error) {
	var cjb bytes.Buffer
	clusterJinjaTmpl, err := template.New("cluster.jinja").Parse(clusterJinja)
	if err != nil {
		return nil, err
	}
	err = clusterJinjaTmpl.Execute(&cjb, s)
	if err != nil {
		return nil, err
	}

	var cyb bytes.Buffer
	clusterYamlTmpl, err := template.New("cluster.yaml").Parse(clusterYaml)
	if err != nil {
		return nil, err
	}
	err = clusterYamlTmpl.Execute(&cyb, s)
	if err != nil {
		return nil, err
	}

	var pdb bytes.Buffer
	if len(s.Volumes) != 0 {
		volumesTmpl, err := template.New("volumes.yaml").Funcs(funcMap).Parse(persistentVolumeJinja)
		if err != nil {
			return nil, err
		}

		err = volumesTmpl.Execute(&pdb, s)
		if err != nil {
			return nil, err
		}
	}

	return map[string][]byte{
		"cluster.jinja": cjb.Bytes(),
		"cluster.yaml":  cyb.Bytes(),
		"volumes.yaml":  pdb.Bytes(),
	}, nil
}
