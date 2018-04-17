package gke

import (
	"bytes"
	"html/template"

	"github.com/hasura/kubeformation/pkg/provider"
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

// TODO: Need to have a default config method

func (s *Spec) GetType() provider.ProviderType {
	return provider.GKE
}

func (s *Spec) MarshalYaml() (map[string][]byte, error) {
	var cjb bytes.Buffer
	clusterJinjaTmpl := template.Must(template.New("cluster.jinja").Parse(clusterJinja))
	err := clusterJinjaTmpl.Execute(&cjb, s)
	if err != nil {
		return nil, err
	}

	var cyb bytes.Buffer
	clusterYamlTmpl := template.Must(template.New("cluster.yaml").Parse(clusterYaml))
	err = clusterYamlTmpl.Execute(&cyb, s)
	if err != nil {
		return nil, err
	}

	return map[string][]byte{
		"cluster.jinja": cjb.Bytes(),
		"cluster.yaml":  cyb.Bytes(),
	}, nil
}
