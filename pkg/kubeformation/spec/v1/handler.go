package v1

import (
	"github.com/ghodss/yaml"
	"github.com/hasura/kubeformation/pkg/kubeformation/spec"
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
		log.WithError(err).Errorf("registering %s version handler failed")
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
