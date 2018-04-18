package provider

// ProviderType denotes each managed kubernetes provider
type ProviderType int

const (
	// NOP indicates that it is not a valid provider.
	NOP ProviderType = iota

	// GKE indicates Google Kubernetes Engine. The template generated will
	// be Google Deployment Manager Template.
	GKE

	// AKS denotes Azure Kubernetes Service. Template generated will be
	// Azure Resource Manger Templates.
	AKS

	// EKS denotes Amazon Elastic Kubernetes Service. Template generated will
	// be AWS CloudFormation Templates.
	EKS
)

type Provider interface {
	GetType() ProviderType
	MarshalFiles() (map[string][]byte, error)
}
