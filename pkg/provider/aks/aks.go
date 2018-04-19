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

var funcMap = template.FuncMap{
	"sub": func(i int) int {
		if i == 0 {
			return 0
		}
		return i - 1
	},
}

type Spec struct {
	Name       string
	K8SVersion string
	NodePools  []NodePool
	Volumes    []Volume
}

type NodePool struct {
	Name   string
	Count  int64
	VMSize string
	OSType string
}

type Volume struct {
	Name   string
	SizeGB int
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
		Volumes: []Volume{},
	}
}

func (s *Spec) GetType() provider.ProviderType {
	return provider.AKS
}

func (s *Spec) MarshalFiles() (map[string][]byte, error) {
	var adb bytes.Buffer
	azureDeployJinjaTmpl, err := template.New("azureDeploy.json").Funcs(funcMap).Parse(azureDeployJinja)
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

	//FIXME: Create new template for managed disks.

	return map[string][]byte{
		"azuredeploy.json":            adb.Bytes(),
		"azuredeploy.parameters.json": pb.Bytes(),
	}, nil
}
