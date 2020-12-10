package torrentleech

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
	"regexp"
	"strings"
	"time"
)

const driverName = "torrentleech"

var (
	rxParse *regexp.Regexp
)

type Driver struct {
	cfg    config.TrackerConfig
	client *http.Client
}

func (p Driver) Name() string {
	return driverName
}

func (p Driver) ParseMessage(message string) (*parser.Result, error) {
	match := rxParse.FindStringSubmatch(message)
	if len(match) != 5 {
		return nil, parser.ErrCannotParse
	}
	res := parser.NewResult(p.Name())
	res.Category = parser.FindCategory(match[1], nil)
	res.Tags = append(res.Tags, match[2])
	release := releaseparser.Parse(match[3])
	res.Release = match[3]
	res.Name = release.Title
	res.Episode = release.Episode
	res.Season = release.Season
	res.LinkSite = match[4]
	pcs := strings.Split(res.LinkSite, "/")
	if len(pcs) != 5 {
		return nil, parser.ErrCannotParse
	}
	tID := transform.ToInt(pcs[len(pcs)-1])
	res.LinkDL = fmt.Sprintf("https://www.torrentleech.org/rss/download/%d/%s/%s.torrent",
		tID, p.cfg.Auth, strings.ReplaceAll(match[3], " ", "."))
	res.Group = release.Group
	res.Resolution = release.Resolution
	res.Source = release.SourceGroup
	res.Codec = release.CodecGroup
	return res, nil
}

func (p Driver) Download(result *parser.Result) (*metainfo.MetaInfo, error) {
	return transport.FetchTorrent(p.client, result.LinkDL)
}

func (p Driver) Login() error {
	return nil
}

func New(trackerConfig config.TrackerConfig) (tracker.Driver, error) {
	return &Driver{
		cfg: trackerConfig,
		client: &http.Client{
			Timeout: time.Second * 10,
		}}, nil
}

type initializer struct{}

func (i initializer) New(trackerConfig config.TrackerConfig) (tracker.Driver, error) {
	return New(trackerConfig)
}

func init() {
	rxParse = regexp.MustCompile(`^New Torrent Announcement:\s<(?P<cat>.+?)\s+::\s+(?P<tag>.+?)>\s+Name:'(?P<release>.+?)'.+?\s(\S+)$`)
	tracker.Register(driverName, initializer{})
}
