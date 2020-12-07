package internal

import (
	"context"
	"github.com/leighmacdonald/lurkr/internal/config"
	"github.com/leighmacdonald/lurkr/internal/sources/irc"
	"github.com/leighmacdonald/lurkr/internal/transport"
	filesystem_transport "github.com/leighmacdonald/lurkr/internal/transport/filesystem"
	"github.com/leighmacdonald/lurkr/internal/transport/sftp"
	"github.com/leighmacdonald/lurkr/pkg/filesystem"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	sourceIRC    map[string]*irc.Connection
	sourcesIRCMU *sync.RWMutex
)

func Start(ctx context.Context) {
	for _, cfg := range config.Sources.IRC {
		c := irc.New(ctx, cfg)
		go func() {
			if err := c.Start(); err != nil {
				log.Errorf(err.Error())
			}
		}()
		sourcesIRCMU.Lock()
		sourceIRC[cfg.Name] = c
		sourcesIRCMU.Unlock()
	}
	for _, cfg := range config.Watch {
		go filesystem.WatchDir(ctx, cfg.Path, func(path string) error {
			// TODO needed for race condition w/OS / file copy?
			time.Sleep(2 * time.Second)
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			fileName := filepath.Base(path)
			return SendTransport(cfg.TransportType, cfg.TransportName, f, fileName)
		})
		log.Debugf("Watching files under %s", cfg.Path)
	}
}

func Stop() {
	sourcesIRCMU.RLock()
	defer sourcesIRCMU.RUnlock()
	for _, conn := range sourceIRC {
		conn.Stop()
	}
}

func SendTransport(transportType config.TransportType, transportName string, reader io.Reader, path string) error {
	var tx transport.Transport
	switch transportType {
	case config.SFTP:
		c, err := config.TransportConfigSFTP(transportName)
		if err != nil {
			return errors.Wrapf(transport.ErrInvalidTransport, "Invalid transport config format: %v", err)
		}
		if c == nil {
			return errors.Wrapf(transport.ErrInvalidTransport, "Invalid transport name: %s", transportName)
		}
		trans, err2 := sftp.NewSFTPTransport(c)
		if err2 != nil {
			return err2
		}
		tx = trans
		path = strings.Join([]string{c.Path, path}, "/")
	case config.Filesystem:
		c, err := config.TransportConfigFile(transportName)
		if err != nil {
			return errors.Wrapf(transport.ErrInvalidTransport, "Invalid transport config format: %v", err)
		}
		if c == nil {
			return errors.Wrapf(transport.ErrInvalidTransport, "Invalid transport name: %s", transportName)
		}
		trans, err2 := filesystem_transport.NewFileTransport(c)
		if err2 != nil {
			return err2
		}
		tx = trans
		path = strings.Join([]string{c.Path, path}, "/")
	default:
		return transport.ErrInvalidTransport
	}
	log.Debugf("Moving file to: %s", path)
	return tx.Send(reader, path)
}

func init() {
	sourceIRC = make(map[string]*irc.Connection)
	sourcesIRCMU = &sync.RWMutex{}
}
