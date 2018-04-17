package provider

type AKS struct {
	ProviderType `json:"provider"`
	Version      string `json:"version"`
}

func (p *AKS) GetType() ProviderType {
	return p.ProviderType
}

func (p *AKS) MarshalYaml() (map[string][]byte, error) {
	return nil, nil
}
