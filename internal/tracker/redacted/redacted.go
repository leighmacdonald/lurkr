package redacted

import (
	"github.com/leighmacdonald/lurkr/internal/config"
	"github.com/leighmacdonald/lurkr/internal/parser"
	"github.com/leighmacdonald/lurkr/internal/torrent"
	"github.com/leighmacdonald/lurkr/internal/tracker"
	"github.com/leighmacdonald/lurkr/pkg/transform"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const driverName = "redacted"

type RED struct {
	rxRelease *regexp.Regexp
	client    *http.Client
	baseURL   string
	cfg       *config.TrackerConfig
}

func (p RED) Name() string {
	return driverName
}

func (p RED) Login() error {
	// Uses API Token
	return nil
}

func (p RED) Download(result *parser.Result) (*torrent.File, error) {
	u, err := url.Parse(p.baseURL)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("action", "download")
	q.Set("id", result.LinkDL)
	u.RawQuery = q.Encode()
	log.Debugf("Downloading: %s", u.String())
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if p.cfg.Auth == "" {
		return nil, config.ErrInvalidAuth
	}
	req.Header.Set("Authorization", p.cfg.Auth)
	resp, err2 := p.client.Do(req)
	if err2 != nil {
		return nil, err2
	}
	var tf torrent.File
	if err3 := torrent.Decode(resp.Body, &tf); err3 != nil {
		return nil, errors.Wrapf(err3, "Failed to decode downloaded torrent")
	}
	defer resp.Body.Close()
	return &tf, nil
}

func (p RED) ParseMessage(message string) (*parser.Result, error) {
	result := &parser.Result{Tracker: p.cfg.Name}
	args := strings.Split(message, " - ")
	if len(args) != 5 {
		return nil, parser.ErrCannotParse
	}
	urls := transform.TrimStrings(strings.Split(args[3], " / "))
	if len(urls) != 2 {
		return nil, parser.ErrCannotParse
	}
	result.LinkSite = urls[0]
	u, err := url.ParseRequestURI(result.LinkSite)
	if err != nil {
		return nil, parser.ErrCannotParse
	}
	tId := u.Query().Get("id")
	if tId == "" {
		return nil, parser.ErrCannotParse
	}
	result.LinkDL = tId
	result.Formats = transform.NormalizeStrings(strings.Split(args[2], "/"))
	result.Tags = transform.NormalizeStrings(strings.Split(args[4], ","))
	result.Name = args[0]
	match := p.rxRelease.FindStringSubmatch(args[1])
	if len(match) != 4 {
		return nil, parser.ErrCannotParse
	}
	result.SubName = match[1]
	result.Year = transform.ToInt(match[2])
	return result, nil
}

func New(trackerConfig *config.TrackerConfig) (tracker.Driver, error) {
	return &RED{
		baseURL: "https://redacted.ch/ajax.php",
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		rxRelease: regexp.MustCompile(`(?P<name>.+?)\s+\[(?P<year>\d+)]\s\[(?P<type>.+?)]$`),
	}, nil
}

type initializer struct{}

func (i initializer) New(trackerConfig *config.TrackerConfig) (tracker.Driver, error) {
	return New(trackerConfig)
}

func init() {
	tracker.Register(driverName, initializer{})
}
