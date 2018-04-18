// Package gke implements GKE provider for Kubeformation. It can generate Google
// Deployment Manager template in Jinja format. Only input required is the Spec
// object.
//
// Google Deployment Manager template supports Python files, Jinja templates as
// well as plain Yaml files.
package gke

import (
	"bytes"
	"html/template"

	"github.com/hasura/kubeformation/pkg/provider"
)

const (
	DefaultK8SVersion  = "1.9"
	DefaultMachineType = "n1-standard-1"
	DefaultImageType   = "cos"
)

type Spec struct {
	Name       string
	K8SVersion string
	NodePools  []NodePool
}
type NodePool struct {
	Name        string
	Labels      map[string]string
	Size        int64
	MachineType string
	ImageType   string
}

func NewDefaultSpec() *Spec {
	return &Spec{
		Name:       "gke-cluster",
		K8SVersion: DefaultK8SVersion,
		NodePools: []NodePool{
			NodePool{
				Name:        "gke-cluster-np-1",
				Size:        1,
				MachineType: DefaultMachineType,
				ImageType:   DefaultImageType,
			},
		},
	}
}

func (s *Spec) GetType() provider.ProviderType {
	return provider.GKE
}

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

	return map[string][]byte{
		"cluster.jinja": cjb.Bytes(),
		"cluster.yaml":  cyb.Bytes(),
	}, nil
}
