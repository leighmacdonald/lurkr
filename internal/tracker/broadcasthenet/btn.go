package broadcasthenet

import (
	"fmt"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/cytec/releaseparser"
	"github.com/leighmacdonald/lurkr/internal/config"
	"github.com/leighmacdonald/lurkr/internal/parser"
	"github.com/leighmacdonald/lurkr/internal/tracker"
	"github.com/leighmacdonald/lurkr/internal/transport"
	"github.com/leighmacdonald/lurkr/pkg/transform"
	"net/http"
	"strings"
	"time"
)

const driverName = "broadcasthenet"

type Driver struct {
	cfg    *config.TrackerConfig
	client *http.Client
}

func (p Driver) Name() string {
	return driverName
}

func (p Driver) ParseMessage(message string) (*parser.Result, error) {
	result := parser.NewResult(p.Name())
	args := strings.Split(message, " | ")
	if len(args) != 14 {
		return nil, parser.ErrCannotParse
	}
	parsed := releaseparser.Parse(args[len(args)-1])
	tID := transform.ToInt(args[10])
	result.Name = parsed.Title
	result.Episode = parsed.Episode
	result.Season = parsed.Season
	// We have no way of really getting the id= parameter w/o hitting the api
	result.LinkSite = fmt.Sprintf("https://broadcasthe.net/torrents.php?torrentid=%d", tID)
	result.LinkDL = fmt.Sprintf(
		"https://broadcasthe.net/torrents.php?action=download&id=%d&authkey=%s&torrent_pass=%s",
		tID, p.cfg.Auth, p.cfg.Passkey)
	result.Group = parsed.Group
	result.Category = parser.TV
	result.Year = transform.ToInt(args[3])
	result.Pack = args[2] != "Episode"
	result.Codec = parsed.CodecGroup
	result.Source = parsed.SourceGroup
	result.Tags = []string{}
	return result, nil
}

func (p Driver) Download(result *parser.Result) (*metainfo.MetaInfo, error) {
	return transport.FetchTorrent(p.client, result.LinkDL)
}

func (p Driver) Login() error {
	return nil
}

func New(trackerConfig *config.TrackerConfig) (tracker.Driver, error) {
	return &Driver{
		cfg:    trackerConfig,
		client: &http.Client{Timeout: time.Second * 10},
	}, nil
}

type initializer struct{}

func (i initializer) New(trackerConfig *config.TrackerConfig) (tracker.Driver, error) {
	return New(trackerConfig)
}

func init() {
	tracker.Register(driverName, initializer{})
}
