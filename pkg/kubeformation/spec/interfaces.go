package spec

import "github.com/hasura/kubeformation/pkg/kubeformation/provider"

// VersionedSpecHandler is the interface every handler should implement.
type VersionedSpecHandler interface {
	// Read takes data and gives a VersionedSpecHandler
	Read(data []byte) (VersionedSpecHandler, error)

	// GetVersion returns the current handler version
	GetVersion() string

	// GenerateProviderTemplate returns a map of template files for
	// the given provider. If provider is not explicitly passed,
	// it is taken from the Spec.
	GenerateProviderTemplate(provider.ProviderType) (map[string][]byte, error)
}
