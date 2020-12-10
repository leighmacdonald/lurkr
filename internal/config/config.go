package config

import (
	"github.com/dustin/go-humanize"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

var (
	ErrInvalidConfig = errors.New("Invalid config")
	ErrInvalidAuth   = errors.New("Invalid auth value")
)

type TransportType string

const (
	Filesystem TransportType = "filesystem"
	SFTP       TransportType = "sftp"
)

type rootConfig struct {
	General       GeneralConfig               `mapstructure:"general"`
	Database      DatabaseConfig              `mapstrucutre:"database"`
	Logs          LogConfig                   `mapstructure:"log"`
	Trackers      []TrackerConfig             `mapstructure:"sources"`
	Watch         map[string]WatchConfig      `mapstructure:"watch"`
	TransportSFTP map[string]SFTPConfig       `mapstructure:"transport_sftp"`
	TransportFile map[string]FileSystemConfig `mapstructure:"transport_file"`
}

type GeneralConfig struct {
	Debug bool `mapstructure:"debug"`
	Dry   bool `mapstructure:"dry_run"`
}

type LogConfig struct {
	Level          string `mapstructure:"level"`
	ForceColours   bool   `mapstructure:"force_colours"`
	DisableColours bool   `mapstructure:"disable_colours"`
	ReportCaller   bool   `mapstructure:"report_caller"`
	FullTimestamp  bool   `mapstructure:"full_timestamp"`
}

type FileSystemConfig struct {
	Path string `mapstructure:"path"`
}

type SFTPConfig struct {
	Address  string `mapstructure:"address"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Key      string `mapstructure:"key"`
	Path     string `mapstructure:"path"`
}

type WatchConfig struct {
	Path          string        `mapstructure:"path"`
	TransportName string        `mapstructure:"transport_name"`
	TransportType TransportType `mapstructure:"transport_type"`
}

type IRCConfig struct {
	Enabled bool `mapstructure:"enabled"`
	// Nick is your IRC nick, usually the same as your username
	Nick string `mapstructure:"nick"`
	// Password is your **IRC** password
	Password string `mapstructure:"password"`
	// Perform is a list of commands to run after connecting to IRC such as nickserv identification, Invite bots, etc.
	Perform []string `mapstructure:"perform"`
	Address string   `mapstructure:"address"`
	// SSL enables SSL IRC Connections
	SSL bool `mapstructure:"ssl"`
	// SSLVerify will enforce valid SSL certificates on the IRC server connection
	SSLVerify bool `mapstructure:"ssl_verify"`
	// Channels to listen to announce events eg:  ["#announce", "#tv-announce"]
	Channels []string `mapstructure:"channels"`
	// BotName is the nick of the announce bot
	BotName string `mapstructure:"bot_name"`
	// BotWho is the whois info for the announce bot
	BotWho string `mapstructure:"bot_who"`
}

type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

type RSSConfig struct {
	Enabled bool     `mapstructure:"enabled"`
	Feeds   []string `mapstructure:"feeds"`
}

type TransportConfig struct {
	Name string        `mapstructure:"name"`
	Type TransportType `mapstructure:"type"`
}

type APIConfig struct {
	TMDB struct {
		Enabled bool   `mapstructure:"enabled"`
		Key     string `mapstructure:"key"`
	} `mapstructure:"tmdb"`
}

type FilterConfig struct {
	TagsAllowed         []string `mapstructure:"tags_allowed"`
	TagsExcluded        []string `mapstructure:"tags_excluded"`
	MinSizeStr          string   `mapstructure:"min_size"`
	MaxSizeStr          string   `mapstructure:"max_size"`
	MinSize             uint64
	MaxSize             uint64
	FormatsAllowed      []string `mapstructure:"formats_allowed"`
	FormatsExcluded     []string `mapstructure:"formats_excluded"`
	ResolutionsAllowed  []string `mapstructure:"resolutions_allowed"`
	ResolutionsExcluded []string `mapstructure:"resolutions_excluded"`
	CategoriesAllowed   []string `mapstructure:"categories_allowed"`
	CategoriesExcluded  []string `mapstructure:"categories_excluded"`
	EpisodesAllowed     bool     `mapstructure:"episodes_allowed"`
	SeasonsAllowed      bool     `mapstructure:"seasons_allowed"`
	TMDBMinScore        float64  `mapstructure:"tmdb_min_score"`
}

type AuthConfig struct {
	// Username is your site username
	Username string `mapstructure:"username"`
	// AuthToken can be several things depending on the tracker:
	// username:password - for sites that dont have a real API, use your standard username and password when logging in
	// api_key - For sites that provides a API for downloading
	AuthToken string `mapstructure:"auth"`
	// Passkey is your torrent passkey, some sites require this when downloading the .torrent
	Passkey string `mapstructure:"passkey"`
}

type TrackerConfig struct {
	Name      string          `mapstructure:"name"`
	Auth      AuthConfig      `mapstructure:"auth"`
	IRC       IRCConfig       `mapstructure:"irc"`
	RSS       RSSConfig       `mapstructure:"rss"`
	Transport TransportConfig `mapstructure:"transport"`
	Filters   FilterConfig    `mapstructure:"filters"`
}

func Tracker(tracker string) (TrackerConfig, error) {
	for _, t := range Trackers {
		if strings.ToLower(t.Name) == strings.ToLower(tracker) {
			return t, nil
		}
	}
	return TrackerConfig{}, ErrInvalidConfig
}

func TransportConfigSFTP(transportName string) (SFTPConfig, error) {
	v, ok := TransportSFTP[transportName]
	if !ok {
		return SFTPConfig{}, errors.Wrapf(ErrInvalidConfig, "Invalid transport key: %v", transportName)
	}
	return v, nil
}

func TransportConfigFile(transportName string) (FileSystemConfig, error) {
	v, ok := TransportFile[transportName]
	if !ok {
		return FileSystemConfig{}, errors.Wrapf(ErrInvalidConfig, "Invalid transport key: %v", transportName)
	}
	return v, nil
}

var (
	// Default, exported configs
	General = GeneralConfig{
		Debug: false,
		Dry:   true,
	}
	Database = DatabaseConfig{DSN: "lurkr.db"}
	Logs     = LogConfig{
		Level:          "info",
		ForceColours:   false,
		DisableColours: false,
		ReportCaller:   false,
		FullTimestamp:  false,
	}
	Trackers []TrackerConfig
	Watch    = map[string]WatchConfig{}

	TransportSFTP = map[string]SFTPConfig{}
	TransportFile = map[string]FileSystemConfig{}
)

// Read reads in config file and ENV variables if set.
func Read(cfgFile string) error {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatalf("Failed to get HOME dir: %v", err)
		}
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName("lurkr")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrapf(err, "Could not read config")
	}
	var cfg rootConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return errors.Wrapf(err, "Invalid config file format: %v", err)
	}
	if validationErr := validate(cfg); validationErr != nil {
		log.Fatalf("Config failed validation: %v", validationErr)
	}
	General = cfg.General
	Database = cfg.Database
	Logs = cfg.Logs
	Trackers = cfg.Trackers
	TransportSFTP = cfg.TransportSFTP
	TransportFile = cfg.TransportFile
	Watch = cfg.Watch
	configureLogger(log.StandardLogger())
	log.Infof("Using config file: %s", viper.ConfigFileUsed())
	return nil
}

func validate(cfg rootConfig) error {
	for _, t := range cfg.Trackers {
		if t.Filters.MaxSizeStr != "" {
			s, e := humanize.ParseBytes(t.Filters.MaxSizeStr)
			if e != nil {
				return errors.Errorf("Failed to parse max_size for %s: %v", t.Name, e)
			}
			t.Filters.MaxSize = s
		}
		if t.Filters.MinSizeStr != "" {
			s, e := humanize.ParseBytes(t.Filters.MinSizeStr)
			if e != nil {
				return errors.Errorf("Failed to parse min_size for %s: %v", t.Name, e)
			}
			t.Filters.MinSize = s
		}
	}
	return nil
}

func configureLogger(l *log.Logger) {
	level, err := log.ParseLevel(Logs.Level)
	if err != nil {
		log.Fatalf("Invalid log level: %s", Logs.Level)
	}
	l.SetLevel(level)
	l.SetFormatter(&log.TextFormatter{
		ForceColors:   Logs.ForceColours,
		DisableColors: Logs.DisableColours,
		FullTimestamp: Logs.FullTimestamp,
	})
	l.SetReportCaller(Logs.ReportCaller)
}
