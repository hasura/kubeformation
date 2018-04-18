package aks

import "github.com/hasura/kubeformation/pkg/provider"

type Spec struct {
}

func NewDefaultSpec() *Spec {
	return &Spec{}
}

func (s *Spec) GetType() provider.ProviderType {
	return provider.AKS
}

func (s *Spec) MarshalFiles() (map[string][]byte, error) {
	return nil, nil
}
