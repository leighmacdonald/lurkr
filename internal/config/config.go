package config

import (
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
	Logs          LogConfig                   `mapstructure:"log"`
	Source        SourcesConfig               `mapstructure:"sources"`
	Watch         map[string]WatchConfig      `mapstructure:"watch"`
	TransportSFTP map[string]SFTPConfig       `mapstructure:"transport_sftp"`
	TransportFile map[string]FileSystemConfig `mapstructure:"transport_file"`
}

type SourcesConfig struct {
	IRC map[string]TrackerConfig `mapstructure:"irc"`
	RSS map[string]TrackerConfig `mapstructure:"rss"`
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

type TrackerConfig struct {
	// Name of the tracker, should match the driverName of the tracker
	Name string `mapstructure:"name"`
	// Nick is your IRC nick, usually the same as your username
	Nick string `mapstructure:"nick"`
	// Username is your site username
	Username string `mapstructure:"username"`
	// Password is your **IRC** password
	Password string `mapstructure:"password"`
	// Perform is a list of commands to run after connecting to IRC such as Nickserv identification, Invite bots, etc.
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
	// Auth can be several things depending on the tracker:
	// username:password - for sites that dont have a real API, use your standard username and password when logging in
	// api_key - For sites that provides a API for downloading
	Auth string `mapstructure:"auth"`
	// Passkey is your torrent passkey, some sites require this when downloading the .torrent
	Passkey string `mapstructure:"passkey"`
	// TransportName should refer to one of the transport names defined in the config file
	TransportName string `mapstructure:"transport_name"`
	// TransportType should refer to one of the transport types defined in the config file
	TransportType string   `mapstructure:"transport_type"`
	RSSFeeds      []string `mapstructure:"rss_feeds"`
}

func Tracker(tracker string) (TrackerConfig, error) {
	for _, t := range Sources.IRC {
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
	Logs = LogConfig{
		Level:          "info",
		ForceColours:   false,
		DisableColours: false,
		ReportCaller:   false,
		FullTimestamp:  false,
	}
	Sources = SourcesConfig{
		IRC: nil,
		RSS: nil,
	}
	Watch = map[string]WatchConfig{}

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
	General = cfg.General
	Logs = cfg.Logs
	Sources = cfg.Source
	TransportSFTP = cfg.TransportSFTP
	TransportFile = cfg.TransportFile
	Watch = cfg.Watch
	configureLogger(log.StandardLogger())
	log.Infof("Using config file: %s", viper.ConfigFileUsed())
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
