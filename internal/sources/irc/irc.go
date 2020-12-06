package irc

import (
	"context"
	"crypto/tls"
	"github.com/leighmacdonald/lurkr/internal/config"
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
	log.Debugln(event)
}

func on464(event *irc.Event, cfg config.TrackerConfig) {
	event.Connection.SendRawf("PASS %s", cfg.Password)
}

func onJOIN(event *irc.Event, cfg config.TrackerConfig) {
	log.Infof("Joined channel: %s", event.MessageWithoutFormat())
}

func onPRIVMSG(event *irc.Event, cfg config.TrackerConfig) {
	valid := false
	for _, allowedChannel := range cfg.Channels {
		if strings.ToLower(event.Arguments[0]) == strings.ToLower(allowedChannel) {
			valid = true
			break
		}
	}
	if valid {
		log.Infof("[%s-%s] %s: %s", cfg.Name, event.Arguments[0], event.Nick, event.MessageWithoutFormat())
	}
}

func onQUIT(event *irc.Event, cfg config.TrackerConfig) {

}

func ignoredEvent(_ *irc.Event) {}

func New(ctx context.Context, cfg config.TrackerConfig) *Connection {
	conn := irc.IRC(cfg.Username, cfg.Username)
	conn.UseTLS = cfg.SSL
	conn.TLSConfig = &tls.Config{InsecureSkipVerify: !cfg.SSLVerify}
	conn.VerboseCallbackHandler = config.General.Debug
	conn.Debug = config.General.Debug
	conn.Log = log.StandardLogger()
	ignoredEvents := []string{
		"001", "002", "003", "004", "005",
		"251", "252", "254", "255", "265", "266",
		"332", "333", "353", "366", "375", "376",
		"MODE", "PING", "PONG", "CTCP_ACTION",
	}
	knownEvents := map[string]func(e *irc.Event){
		"336":     func(e *irc.Event) { on336(e, cfg) },
		"464":     func(e *irc.Event) { on464(e, cfg) },
		"JOIN":    func(e *irc.Event) { onJOIN(e, cfg) },
		"QUIT":    func(e *irc.Event) { onQUIT(e, cfg) },
		"PRIVMSG": func(e *irc.Event) { onPRIVMSG(e, cfg) },
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
			log.Debugf("Unhandled event (%s): %s", event.Code, event.MessageWithoutFormat())
		}
	})
	return &Connection{
		name: cfg.Name,
		conn: conn,
		ctx:  ctx,
		cfg:  cfg,
	}
}

func (i Connection) Stop() {
	i.conn.Disconnect()
}

func (i Connection) Start() error {
	if err := i.conn.Connect(i.cfg.Address); err != nil {
		return errors.Wrapf(err, "Failed to connect to server: %v", err)
	}
	i.conn.Loop()
	return nil
}
