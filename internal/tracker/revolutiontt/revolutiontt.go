package revolutiontt

import (
	"fmt"
	"github.com/cytec/releaseparser"
	"github.com/leighmacdonald/lurkr/internal/config"
	"github.com/leighmacdonald/lurkr/internal/parser"
	"github.com/leighmacdonald/lurkr/internal/torrent"
	"github.com/leighmacdonald/lurkr/internal/tracker"
	"github.com/leighmacdonald/lurkr/pkg/transform"
	"github.com/pkg/errors"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const driverName = "revolutiontt"

type RevTT struct {
	baseURL  string
	username string
	password string
	cfg      *config.TrackerConfig
	client   *http.Client
}

func (p RevTT) Name() string {
	return driverName
}

func (p RevTT) Download(result *parser.Result) (*torrent.File, error) {
	req, err := http.NewRequest("GET", result.LinkSite, nil)
	if err != nil {
		return nil, err
	}
	resp, err2 := p.client.Do(req)
	if err2 != nil {
		return nil, err2
	}
	_, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		return nil, err3
	}
	defer resp.Body.Close()
	return nil, nil
}

func (p RevTT) ParseMessage(message string) (*parser.Result, error) {
	result := &parser.Result{}
	args := strings.Split(message, " | ")
	result.Tags = transform.NormalizeStrings(strings.Split(args[1], "/"))
	result.LinkSite = args[2]
	parsed := releaseparser.Parse(strings.Replace(args[0], "!new ", "", -1))
	result.Name = parsed.Title
	result.Year = parsed.Year
	return nil, parser.ErrCannotParse
}

func (p RevTT) Login() error {
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

func New(cfg *config.TrackerConfig) (*RevTT, error) {
	userInfo := strings.SplitN(cfg.Auth, ":", 2)
	if len(userInfo) != 2 {
		return nil, errors.New("Invalid credentials")
	}
	// Store cookie
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create cookiejar")
	}
	return &RevTT{
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
