// Package provider implements the logic of generating cloud platform specific
// templates.
package provider

// ProviderType indicates a particular managed kubernetes provider.
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

// Provider interface should be implemented by a provider. These methods are
// called by the spec package
type Provider interface {
	// GetType returns the provider type of the implementing provider
	GetType() ProviderType

	// MarshalFiles render the template files required by the provider. It
	// returns a map of file name and data, each entry denoting a template file.
	// Returns an error if the rendering goes wrong.
	MarshalFiles() (map[string][]byte, error)
}

// Parse takes a string p and return the corresponding provider type
func Parse(p string) ProviderType {
	switch p {
	case "gke":
		return GKE
	case "aks":
		return AKS
	case "eks":
		return EKS
	default:
		return NOP
	}
}
