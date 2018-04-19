package v1

import (
	"github.com/ghodss/yaml"
	"github.com/hasura/kubeformation/pkg/provider"
	"github.com/hasura/kubeformation/pkg/provider/aks"
	"github.com/hasura/kubeformation/pkg/provider/gke"
	"github.com/hasura/kubeformation/pkg/spec"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// ErrSpecParseFailure can be thrown if the given spec cannot be parsed into
// a known format
var ErrSpecParseFailure = errors.New("kubeformation: parsing spec failed")

// version string for this handler
const version string = "v1"

// init registers the handler with main module
func init() {
	log.Debugf("registering %s version handler", version)
	s := ClusterSpec{Version: version}
	err := spec.Register(version, &s)
	if err != nil {
		log.WithError(err).Errorf("registering %s version handler failed", version)
	}
}

// GetVersion returns the current spec version
func (s *ClusterSpec) GetVersion() string {
	return s.Version
}

// Read returns a Handler after reading the spec
func (s *ClusterSpec) Read(data []byte) (spec.VersionedSpecHandler, error) {
	err := yaml.Unmarshal(data, s)
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}
	return s, nil
}

// GenerateProviderTemplate returns a map of template files for
// the given provider. If provider is not explicitly passed,
// it is taken from the Spec.
func (s *ClusterSpec) GenerateProviderTemplate(providerType provider.ProviderType) (map[string][]byte, error) {
	if providerType == provider.NOP {
		switch s.Provider {
		case "gke":
			providerType = provider.GKE
		case "aks":
			providerType = provider.AKS
		case "eks":
			providerType = provider.EKS
		}
	}
	switch providerType {
	case provider.GKE:
		spec := gke.NewDefaultSpec()
		spec.Name = s.Name
		spec.K8SVersion = s.K8SVersion
		if len(s.Nodes) > 0 {
			spec.NodePools = []gke.NodePool{}
		}
		for _, nodePool := range s.Nodes {
			pool := gke.NodePool{
				Name:        nodePool.Name,
				MachineType: nodePool.Type,
				ImageType:   nodePool.OSImage,
				Labels:      nodePool.Labels,
				Size:        nodePool.PoolSize,
			}
			spec.NodePools = append(spec.NodePools, pool)
		}
		return spec.MarshalFiles()
	case provider.AKS:
		spec := aks.NewDefaultSpec()
		spec.Name = s.Name
		spec.K8SVersion = s.K8SVersion
		// Add Nodes
		if len(s.Nodes) > 0 {
			spec.NodePools = []aks.NodePool{}
		}
		for _, nodePool := range s.Nodes {
			pool := aks.NodePool{
				Name:   nodePool.Name,
				VMSize: nodePool.Type,
				Count:  nodePool.PoolSize,
				OSType: nodePool.OSImage,
			}
			spec.NodePools = append(spec.NodePools, pool)
		}
		// Add Volumes
		if len(s.Volumes) > 0 {
			spec.Volumes = []aks.Volume{}
		}
		for _, volume := range s.Volumes {
			disk := aks.Volume{
				Name:   volume.Name,
				SizeGB: volume.Size,
			}
			spec.Volumes = append(spec.Volumes, disk)
		}
		return spec.MarshalFiles()
	}
	return nil, errors.New("unknown provider")
}
