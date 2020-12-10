package irc

import (
	"context"
	"crypto/tls"
	"github.com/leighmacdonald/lurkr/internal/config"
	"github.com/leighmacdonald/lurkr/internal/parser"
	"github.com/leighmacdonald/lurkr/internal/tracker"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thoj/go-ircevent"
	"strings"
)

type Connection struct {
	name string
	conn *irc.Connection
	ctx  context.Context
	cfg  config.TrackerConfig
}

func on336(event *irc.Event, cfg config.TrackerConfig) {
	log.Debugf("[%s] %s", cfg.Name, event.MessageWithoutFormat())
}

func on464(event *irc.Event, cfg config.TrackerConfig) {
	event.Connection.SendRawf("PASS %s", cfg.IRC.Password)
}

func onJOIN(event *irc.Event, cfg config.TrackerConfig) {
	if strings.ToLower(event.Nick) == strings.ToLower(cfg.IRC.Nick) {
		log.Infof("[%s] Joined channel: %s", cfg.Name, event.MessageWithoutFormat())
	}
}

func onPRIVMSG(event *irc.Event, cfg config.TrackerConfig, driver tracker.Driver, announceChan chan parser.Announce) {
	valid := false
	for _, allowedChannel := range cfg.IRC.Channels {
		if strings.ToLower(event.Arguments[0]) == strings.ToLower(allowedChannel) {
			valid = true
			break
		}
	}
	if !valid {
		return
	}
	log.Infof("[%s-%s] %s: %s", cfg.Name, event.Arguments[0], event.Nick, event.MessageWithoutFormat())
	result, err := driver.ParseMessage(event.MessageWithoutFormat())
	if err != nil {
		log.Debugf("[%s] Failed to parse release from message: %v", cfg.Name, err)
		return
	}
	announceChan <- parser.Announce{Parsed: result, Cfg: cfg}
}

func onQUIT(event *irc.Event, cfg config.TrackerConfig) {

}

func ignoredEvent(_ *irc.Event) {}

type handlerFunc func(e *irc.Event)

func New(ctx context.Context, cfg config.TrackerConfig, announceChan chan parser.Announce) (*Connection, error) {
	driver, err := tracker.New(cfg)
	if err != nil {
		return nil, err
	}
	conn := irc.IRC(cfg.Auth.Username, cfg.Auth.Username)
	conn.UseTLS = cfg.IRC.SSL
	conn.TLSConfig = &tls.Config{InsecureSkipVerify: !cfg.IRC.SSLVerify}
	conn.VerboseCallbackHandler = config.General.Debug
	conn.Debug = config.General.Debug
	conn.Log = log.StandardLogger()
	ignoredEvents := []string{
		"001", "002", "003", "004", "005",
		"251", "252", "254", "255", "265", "266",
		"332", "333", "353", "366", "372", "375", "376",
		"MODE", "PING", "PONG", "CTCP_ACTION",
	}

	knownEvents := map[string]handlerFunc{
		"336":     func(e *irc.Event) { on336(e, cfg) },
		"464":     func(e *irc.Event) { on464(e, cfg) },
		"JOIN":    func(e *irc.Event) { onJOIN(e, cfg) },
		"QUIT":    func(e *irc.Event) { onQUIT(e, cfg) },
		"PRIVMSG": func(event *irc.Event) { onPRIVMSG(event, cfg, driver, announceChan) },
		"NOTICE": func(e *irc.Event) {
			log.Infof("[%s] NOTICE: %s", cfg.Name, e.MessageWithoutFormat())
		},
	}
	for code, fn := range knownEvents {
		conn.AddCallback(code, fn)
		ignoredEvents = append(ignoredEvents, code)
	}

	conn.AddCallback("*", func(event *irc.Event) {
		known := false
		for _, ie := range ignoredEvents {
			if ie == event.Code {
				known = true
				break
			}
		}
		if !known {
			log.Debugf("[%s] Unhandled event (%s): %s", cfg.Name, event.Code, event.MessageWithoutFormat())
		}
	})
	return &Connection{
		name: cfg.Name,
		conn: conn,
		ctx:  ctx,
		cfg:  cfg,
	}, nil
}

func (i Connection) Stop() {
	i.conn.Disconnect()
}

func (i Connection) Start() error {
	if err := i.conn.Connect(i.cfg.IRC.Address); err != nil {
		return errors.Wrapf(err, "Failed to connect to server: %v", err)
	}
	i.conn.Loop()
	return nil
}
