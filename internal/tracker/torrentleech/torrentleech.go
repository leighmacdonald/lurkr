package torrentleech

import (
	"github.com/leighmacdonald/lurkr/internal/config"
	"github.com/leighmacdonald/lurkr/internal/parser"
	"github.com/leighmacdonald/lurkr/internal/torrent"
	"github.com/leighmacdonald/lurkr/internal/tracker"
)

const driverName = "torrentleech"

type TorrentLeech struct{}

func (p TorrentLeech) Name() string {
	return driverName
}

func (p TorrentLeech) ParseMessage(message string) (*parser.Result, error) {
	panic("implement me")
}

func (p TorrentLeech) Download(result *parser.Result) (*torrent.File, error) {
	panic("implement me")
}

func (p TorrentLeech) Login() error {
	panic("implement me")
}

func New(trackerConfig *config.TrackerConfig) (tracker.Driver, error) {
	return &TorrentLeech{}, nil
}

type initializer struct{}

func (i initializer) New(trackerConfig *config.TrackerConfig) (tracker.Driver, error) {
	return New(trackerConfig)
}

func init() {
	tracker.Register(driverName, initializer{})
}
