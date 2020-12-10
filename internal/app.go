package internal

import (
	"bytes"
	"context"
	"github.com/dustin/go-humanize"
	"github.com/leighmacdonald/lurkr/internal/config"
	"github.com/leighmacdonald/lurkr/internal/parser"
	"github.com/leighmacdonald/lurkr/internal/sources/irc"
	"github.com/leighmacdonald/lurkr/internal/tracker"
	"github.com/leighmacdonald/lurkr/internal/transport"
	filesystemtransport "github.com/leighmacdonald/lurkr/internal/transport/filesystem"
	"github.com/leighmacdonald/lurkr/internal/transport/sftp"
	"github.com/leighmacdonald/lurkr/pkg/filesystem"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
	newAnnounce  chan parser.Announce
	stopChan     chan bool
	db           *gorm.DB
)

func announceWorker() {
	for {
		select {
		case nr := <-newAnnounce:
			log.Debugf("[%s] Got new release: %v", nr.Cfg.Name, nr.Parsed.Release)
			if !parser.Match(nr.Cfg, nr.Parsed) {
				log.Debugf("Skipped release: %v", nr.Parsed.Release)
				continue
			}
			driver, err := tracker.New(nr.Cfg)
			if err != nil {
				log.Fatalf("Failed to create tracker driver for download: %v", err)
			}
			mi, err := driver.Download(nr.Parsed)
			if err != nil {
				log.Errorf("Failed to download release: %v", err)
				continue
			}
			i, err := mi.UnmarshalInfo()
			totalSize := uint64(i.TotalLength())
			if nr.Cfg.Filters.MinSize > 0 && totalSize < nr.Cfg.Filters.MinSize {
				log.Debugf("Skippred release (too small): %s", humanize.Bytes(totalSize))
				continue
			}
			if nr.Cfg.Filters.MaxSize > 0 && totalSize > nr.Cfg.Filters.MaxSize {
				log.Debugf("Skippred release (too large): %s", humanize.Bytes(totalSize))
				continue
			}

			rls, err := GetRelease(mi.HashInfoBytes().HexString())
			if err != nil && err.Error() != "record not found" {
				log.Debugf("Skipped duplicate release: %v", nr.Parsed.Name)
				continue
			}
			if rls.Hash != "" {
				log.Debugf("Release already downloaded: %s", nr.Parsed.Release)
				continue
			}
			b := bytes.NewBuffer(nil)
			if err := mi.Write(b); err != nil {
				log.Errorf("Failed to encode .torrent for transport: %v", err)
				continue
			}
			if err := SendTransport(nr.Cfg.Transport.Type, nr.Cfg.Transport.Name, b, i.Name); err != nil {
				log.Errorf("Failed to send .torrent over transport (%s:%s): %v",
					nr.Cfg.Transport.Type, nr.Cfg.Transport.Name, err)
				continue
			}
			// Save release
			if err := InsertRelease(mi.HashInfoBytes().HexString(), driver.Name()); err != nil {
				log.Errorf("Failed to insert release into database: %v", err)
			}
		case <-stopChan:
			return
		}
	}
}

func Start(ctx context.Context) {
	newDb, err := NewDb(config.Database.DSN)
	if err != nil {
		log.Fatalf("Failed to connecto to database: %v", err)
	}
	db = newDb
	c, cancel := context.WithCancel(ctx)
	for _, cfg := range config.Trackers {
		if !cfg.IRC.Enabled {
			continue
		}
		ircConn, err := irc.New(c, cfg, newAnnounce)
		if err != nil {
			log.Fatalf("Error creating irc conn: %v", err)
		}
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
	go announceWorker()
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
		trans, err2 := filesystemtransport.NewFileTransport(c)
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
	newAnnounce = make(chan parser.Announce)
}
