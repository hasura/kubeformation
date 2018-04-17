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

// Register a version handler vsh that can handler a particular version of spec
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
