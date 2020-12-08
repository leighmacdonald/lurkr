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
	stopChan     chan bool
)

func Start(ctx context.Context) {
	c, cancel := context.WithCancel(ctx)
	for _, cfg := range config.Sources.IRC {
		ircConn := irc.New(c, cfg)
		go func() {
			if err := ircConn.Start(); err != nil {
				log.Errorf(err.Error())
			}
		}()
		sourcesIRCMU.Lock()
		sourceIRC[cfg.Name] = ircConn
		sourcesIRCMU.Unlock()
	}
	for _, cfg := range config.Watch {
		go filesystem.WatchDir(c, cfg.Path, func(path string) error {
			// TODO needed for race condition w/OS / file copy?
			time.Sleep(2 * time.Second)
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer func() {
				if errF := f.Close(); errF != nil {
					log.Errorf("Failed to close file: %v", errF)
				}
			}()
			fileName := filepath.Base(path)
			return SendTransport(cfg.TransportType, cfg.TransportName, f, fileName)
		})
		log.Debugf("Watching files under %s", cfg.Path)
	}
	<-stopChan
	cancel()
}

func Stop() {
	sourcesIRCMU.RLock()
	defer sourcesIRCMU.RUnlock()
	for _, conn := range sourceIRC {
		conn.Stop()
	}
	stopChan <- true
}

func SendTransport(transportType config.TransportType, transportName string, reader io.Reader, path string) error {
	var tx transport.Transport
	switch transportType {
	case config.SFTP:
		c, err := config.TransportConfigSFTP(transportName)
		if err != nil {
			return errors.Wrapf(transport.ErrInvalidTransport, "Invalid transport config format: %v", err)
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
	stopChan = make(chan bool)
}
