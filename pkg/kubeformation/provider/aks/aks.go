package aks

import "github.com/hasura/kubeformation/pkg/kubeformation/provider"

type Spec struct {
}

func (s *Spec) GetType() provider.ProviderType {
	return provider.AKS
}

func (s *Spec) MarshalYaml() (map[string][]byte, error) {
	return nil, nil
}
