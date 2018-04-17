package provider

type ProviderType int

const (
	ProviderGKE ProviderType = iota
	ProviderAKS
)

type Provider interface {
	GetType() ProviderType
	MarshalYaml() (map[string][]byte, error)
}
