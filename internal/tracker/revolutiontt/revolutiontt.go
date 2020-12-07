package revolutiontt

import (
	"fmt"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/cytec/releaseparser"
	"github.com/leighmacdonald/lurkr/internal/config"
	"github.com/leighmacdonald/lurkr/internal/parser"
	"github.com/leighmacdonald/lurkr/internal/tracker"
	"github.com/leighmacdonald/lurkr/internal/transport"
	"github.com/leighmacdonald/lurkr/pkg/transform"
	"github.com/pkg/errors"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const driverName = "revolutiontt"

type Driver struct {
	baseURL  string
	username string
	password string
	cfg      *config.TrackerConfig
	client   *http.Client
}

func (p Driver) Name() string {
	return driverName
}

func (p Driver) Download(result *parser.Result) (*metainfo.MetaInfo, error) {
	return transport.FetchTorrent(p.client, result.LinkDL)
}

func (p Driver) ParseMessage(message string) (*parser.Result, error) {
	result := parser.NewResult(driverName)
	args := strings.Split(message, " | ")
	tags := transform.NormalizeStrings(strings.Split(args[1], "/"))
	if len(tags) != 2 {
		return nil, parser.ErrCannotParse
	}
	result.Category = parser.FindCategory(tags[0], nil)
	result.Tags = tags[1:]
	result.LinkSite = args[2]
	parsed := releaseparser.Parse(strings.Replace(args[0], "!new ", "", -1))
	result.Name = parsed.Title
	if parsed.Year == 0 {
		result.Year = transform.FindYear(parsed.Title)
	} else {
		result.Year = parsed.Year
	}
	u, err := url.Parse(result.LinkSite)
	if err != nil {
		return nil, parser.ErrCannotParse
	}
	tID := u.Query().Get("id")
	result.LinkDL = fmt.Sprintf("https://revolutiontt.me/download.php/%s?passkey=%s", tID, p.cfg.Passkey)
	return result, nil
}

func (p Driver) Login() error {
	// Login
	// Get login page to acquire a session key first
	resp, err := p.client.Get(fmt.Sprintf("%s/login.php", p.baseURL))
	if err != nil || resp.StatusCode != 200 {
		return errors.Wrapf(err, "failed to load login page")
	}
	// Do login form
	params := url.Values{}
	params.Set("username", p.username)
	params.Set("password", p.password)
	params.Set("submit", "login")
	respLogin, err := p.client.PostForm(fmt.Sprintf("%s/takelogin.php", p.baseURL), params)
	if err != nil {
		return errors.Wrapf(err, "failed to post login")
	}
	if respLogin.StatusCode != 200 {
		return errors.New("Invalid login response code")
	}
	return nil
}

func New(cfg *config.TrackerConfig) (*Driver, error) {
	userInfo := strings.SplitN(cfg.Auth, ":", 2)
	if len(userInfo) != 2 {
		return nil, errors.Wrap(config.ErrInvalidConfig, "auth must be configured as username:password")
	}
	// Store cookie
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create cookiejar")
	}
	if cfg.Passkey == "" {
		return nil, errors.Wrapf(config.ErrInvalidConfig, "Passkey must be set")
	}
	return &Driver{
		baseURL:  `https://revolutiontt.me/`,
		username: userInfo[0],
		password: userInfo[1],
		cfg:      cfg,
		client: &http.Client{
			Jar: jar,
		},
	}, nil
}

type initializer struct{}

func (i initializer) New(trackerConfig *config.TrackerConfig) (tracker.Driver, error) {
	return New(trackerConfig)
}

func init() {
	tracker.Register(driverName, initializer{})
}
