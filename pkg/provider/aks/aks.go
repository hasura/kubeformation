package aks

import (
	"bytes"
	"html/template"

	"github.com/hasura/kubeformation/pkg/provider"
)

const (
	DefaultOSType     = "Linux"
	DefaultK8SVersion = "1.8.1"
	DefaultVMSize     = "Standard_D2_v2"
)

type Spec struct {
	Name       string
	K8SVersion string
	NodePools  []NodePool
}

type NodePool struct {
	Name   string
	Count  int64
	VMSize string
	OSType string
}

func NewDefaultSpec() *Spec {
	return &Spec{
		Name:       "aks-cluster",
		K8SVersion: DefaultK8SVersion,
		NodePools: []NodePool{
			NodePool{
				Name:   "aks-cluster-np-1",
				Count:  1,
				OSType: DefaultOSType,
				VMSize: DefaultVMSize,
			},
		},
	}
}

func (s *Spec) GetType() provider.ProviderType {
	return provider.AKS
}

func (s *Spec) MarshalFiles() (map[string][]byte, error) {
	var adb bytes.Buffer
	azureDeployJinjaTmpl, err := template.New("azureDeploy.jinja").Parse(azureDeployJinja)
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

	return map[string][]byte{
		"azuredeploy.json":            adb.Bytes(),
		"azuredeploy.parameters.json": pb.Bytes(),
	}, nil
}
