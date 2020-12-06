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

type rootConfig struct {
	General       GeneralConfig               `mapstructure:"general"`
	Logs          LogConfig                   `mapstructure:"log"`
	Source        SourcesConfig               `mapstructure:"sources"`
	TransportSFTP map[string]SFTPConfig       `mapstructure:"transport_sftp"`
	TransportFile map[string]FileSystemConfig `mapstructure:"transport_file"`
}

type SourcesConfig struct {
	IRC []TrackerConfig `mapstructure:"irc"`
	RSS []TrackerConfig `mapstructure:"rss"`
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
	Dest string `mapstructure:"dest"`
}

type SFTPConfig struct {
	Address  string `mapstructure:"address"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Key      string `mapstructure:"key"`
	Path     string `mapstructure:"path"`
}

type TrackerConfig struct {
	Name          string   `mapstructure:"name"`
	Nick          string   `mapstructure:"nick"`
	Username      string   `mapstructure:"username"`
	Password      string   `mapstructure:"password"`
	Perform       []string `mapstructure:"perform"`
	Address       string   `mapstructure:"address"`
	SSL           bool     `mapstructure:"ssl"`
	SSLVerify     bool     `mapstructure:"ssl_verify"`
	Channels      []string `mapstructure:"channels"`
	BotName       string   `mapstructure:"bot_name"`
	BotWho        string   `mapstructure:"bot_who"`
	Auth          string   `mapstructure:"auth"`
	TransportName string   `mapstructure:"transport_name"`
	TransportType string   `mapstructure:"transport_type"`
}

func Tracker(tracker string) *TrackerConfig {
	for _, t := range Sources.IRC {
		if strings.ToLower(t.Name) == strings.ToLower(tracker) {
			return &t
		}
	}
	return nil
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
	TransportSFTP = map[string]SFTPConfig{}
	TransportFS   = map[string]FileSystemConfig{
		"default": {Dest: "/watch"},
	}
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
	TransportFS = cfg.TransportFile
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
