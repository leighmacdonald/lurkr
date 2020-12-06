package broadcasthenet

import (
	"github.com/leighmacdonald/lurkr/internal/config"
	"github.com/leighmacdonald/lurkr/internal/parser"
	"github.com/leighmacdonald/lurkr/internal/torrent"
	"github.com/leighmacdonald/lurkr/internal/tracker"
)

const driverName = "broadcasthenet"

type BTNDriver struct{}

func (p BTNDriver) Name() string {
	return driverName
}

func (p BTNDriver) ParseMessage(message string) (*parser.Result, error) {
	panic("implement me")
}

func (p BTNDriver) Download(result *parser.Result) (*torrent.File, error) {
	panic("implement me")
}

func (p BTNDriver) Login() error {
	panic("implement me")
}

func New(trackerConfig *config.TrackerConfig) (tracker.Driver, error) {
	return &BTNDriver{}, nil
}

type initializer struct{}

func (i initializer) New(trackerConfig *config.TrackerConfig) (tracker.Driver, error) {
	return New(trackerConfig)
}

func init() {
	tracker.Register(driverName, initializer{})
}
