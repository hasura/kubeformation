package provider

type GKE struct {
	ProviderType `json:"provider"`
	Version      string `json:"version"`
}

func (p *GKE) GetType() ProviderType {
	return p.ProviderType
}

func (p *GKE) MarshalYaml() (map[string][]byte, error) {
	return nil, nil
}
