// Package spec implements methods to parse and process the cluster spec, which
// is the main input to kubeformation. ClusterSpec is versioned. The version is
// identified by
//   version: v1
// in the cluster spec. Each version has a corresponding handler, which
// implements the VersionedSpecHandler interface. Handler should implement
// the logic of parsing the cluster spec from a raw yaml byte array. It should
// also specify how the cluster spec converts to the Spec specific to each
// provider.
//
// Each handler should also register here by invoking the Register function.
// Usually the invocation is present in the init section of handler. The handler
// should then be imported once in code:
//   import (
//     _ "github.com/hasura/kubeformation/pkg/spec/v1"
//   )
//
// Read method of this package can be called with raw yaml data for cluster
// spec. Version is read from the data and the corresponding handler will be
// returned, which can be used to generate the provider specific template.
//
// When cluster spec changes in a way that breaks existing implementation, a new
// handler, say v2 should be created. This way, both v1 and v2 cluster spec can
// he handled by the same tool.
package spec

import (
	"sync"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// ErrVersionAlreadyRegistered is thrown when a version handler that is already registered
// tries to register again
var ErrVersionAlreadyRegistered = errors.New("kubeformation: version is already registered")

// ErrInvalidSpecVersion is thrown when spec with an unknown version is provided
var ErrInvalidSpecVersion = errors.New("kubeformation: invalid spec version")

// handlers is a container to hold all the version handlers
var handlers = make(map[string]VersionedSpecHandler)

// handlerLock can be used to make sure that concurrent modifications to
// handlers map does not happen
var handlerLock sync.Mutex

// Register a version handler vsh that can handle a particular version of spec
func Register(version string, vsh VersionedSpecHandler) error {
	handlerLock.Lock()
	defer handlerLock.Unlock()
	if _, found := handlers[version]; found {
		log.Debugf("version %s already registered", version)
		return errors.Wrap(ErrVersionAlreadyRegistered, version)
	}
	log.Debugf("spec version %s registered", version)
	handlers[version] = vsh
	return nil
}

// Read a spec and return the corresponding version handler
func Read(data []byte) (VersionedSpecHandler, error) {
	var spec VersionedSpec
	log.Debug("Reading version key from spec")
	err := yaml.Unmarshal(data, &spec)
	if err != nil {
		log.Debug(err)
		return nil, errors.Wrap(err, "reading spec failed")
	}
	version := spec.Version
	log.Debugf("version: %s", version)
	if len(version) == 0 {
		return nil, errors.Wrapf(ErrInvalidSpecVersion, "empty version string")
	}
	handlerLock.Lock()
	defer handlerLock.Unlock()
	log.Debug("checking handler map for version handler")
	vh, ok := handlers[version]
	if !ok {
		log.Debug("version not found in handler map")
		return nil, errors.Wrapf(ErrInvalidSpecVersion, "unsupported version: %s", version)
	}
	log.Debug("version found in handler map")
	return vh.Read(data)
}
