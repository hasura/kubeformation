// Package aks implements AKS provider for Kubeformation. It can generate Azure
// Resource Manager templates in JSON format. Spec is the input struct.
//
// Azure Deployment Manager templates are written in JSON formatted files,
// typically a resource definition file and a parameters file. Kubeformation
// converts Spec into azuredeploy.json and azuredeploy.parameters.json files.
// There could be more files depending on the presence of volumes/disks. These
// files are defined as Go templates in template.go and are rendered with Spec
// as the context.
package aks

import (
	"bytes"
	"html/template"

	"github.com/hasura/kubeformation/pkg/provider"
)

const (
	// Default version of Kubernetes to use
	DefaultK8SVersion = "1.8.1"

	// Default Azure VM size
	DefaultVMSize = "Standard_D2_v2"

	// Default OS type, available options are Linux and Windows
	DefaultOSType = "Linux"
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

// Spec defines the context required to render ARM template.
type Spec struct {
	// Name of the cluster
	Name string

	// Kubernetes version for the cluster
	K8SVersion string

	// Node pools to be created
	NodePools []NodePool

	// Persistent disks to be created
	Volumes []Volume
}

// NodePool defines a collection of nodes (VMs) with similar properties that can
// be used by a Kubernetes cluster to schedule workloads.
type NodePool struct {
	// Name of the node pool
	Name string

	// Number of nodes in this pool
	Count int64

	// Size (type) of the VMs
	VMSize string

	// Type of OS to use
	OSType string
}

// Volume defines a Azure Data Disk that can persist data.
type Volume struct {
	// Name of the disk
	Name string

	// Size of the disk in GB
	SizeGB int
}

// NewDefaultSpec returns a spec object which is barely enough to render a valid
// ARM template.
func NewDefaultSpec() *Spec {
	return &Spec{
		Name:       "aks-cluster",
		K8SVersion: DefaultK8SVersion,
		NodePools: []NodePool{
			NodePool{
				Name:   "np-1",
				Count:  1,
				OSType: DefaultOSType,
				VMSize: DefaultVMSize,
			},
		},
	}
}

// GetType returns the type of this provider.
func (s *Spec) GetType() provider.ProviderType {
	return provider.AKS
}

// MarshalFiles returns the rendered ARM template as a map which indicates
// filename as keys and file content as value.
// FIXME: test does not capture the template errors.
func (s *Spec) MarshalFiles() (map[string][]byte, error) {
	var adb bytes.Buffer
	azureDeployJinjaTmpl, err := template.New("azuredeploy.json").Funcs(funcMap).Parse(azureDeployJSON)
	if err != nil {
		return nil, err
	}
	err = azureDeployJinjaTmpl.Execute(&adb, s)
	if err != nil {
		return nil, err
	}

	var pb bytes.Buffer
	parametersJSONTmpl, err := template.New("azuredeploy.parameters.json").Parse(parametersJSON)
	if err != nil {
		return nil, err
	}
	err = parametersJSONTmpl.Execute(&pb, s)
	if err != nil {
		return nil, err
	}

	// FIXME: Create new template for managed disks.

	return map[string][]byte{
		"azuredeploy.json":            adb.Bytes(),
		"azuredeploy.parameters.json": pb.Bytes(),
	}, nil
}
