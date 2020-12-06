package internal

import (
	"context"
	"github.com/leighmacdonald/lurkr/internal/config"
	"github.com/leighmacdonald/lurkr/internal/sources/irc"
	log "github.com/sirupsen/logrus"
	"sync"
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
}

func Stop() {
	sourcesIRCMU.RLock()
	defer sourcesIRCMU.RUnlock()
	for _, conn := range sourceIRC {
		conn.Stop()
	}
}

func init() {
	sourceIRC = make(map[string]*irc.Connection)
	sourcesIRCMU = &sync.RWMutex{}
}
