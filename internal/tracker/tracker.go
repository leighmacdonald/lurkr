package tracker

import (
	"github.com/anacrolix/torrent/metainfo"
	"github.com/leighmacdonald/lurkr/internal/config"
	"github.com/leighmacdonald/lurkr/internal/parser"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"sync"
)

// Driver provides a generic interface for interacting with a tracker
type Driver interface {
	// ParseMessage parses an IRC announce message
	ParseMessage(message string) (*parser.Result, error)

	// Download the raw .torrent file to send to the configured transport
	Download(result *parser.Result) (*metainfo.MetaInfo, error)

	// Login to the site, if required
	Login() error

	Name() string
}

// Initializer will create the tracker driver instance with the provided config
type Initializer interface {
	New(*config.TrackerConfig) (Driver, error)
}

var (
	drivers   map[string]Initializer
	driversMu *sync.RWMutex

	ErrDownload = errors.New("Failed to download torrent")
)

// Register will register a driver as usable
func Register(name string, initializer Initializer) {
	driversMu.Lock()
	defer driversMu.Unlock()
	drivers[name] = initializer
	log.Debugf("Registered driver: %s", name)
}

// New will instantiate the matching underlying drivers if it exists
func New(trackerConfig *config.TrackerConfig) (Driver, error) {
	driversMu.RLock()
	defer driversMu.RUnlock()
	driver, found := drivers[trackerConfig.Name]
	if !found {
		return nil, errors.New("Unregistered driver")
	}
	return driver.New(trackerConfig)
}

func init() {
	drivers = make(map[string]Initializer)
	driversMu = &sync.RWMutex{}
}
