package filesystem

import (
	"context"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	// Batch write events to avoid issues with certain applications / filesystems writing multiple
	// events when saving a file
	// Will wait N seconds before actually writing the event out to the WatchFunc
	batch   map[string]time.Time
	batchMu *sync.RWMutex
)

const (
	minBatchInterval = time.Second * 3
)

type WatchFunc func(path string) error

func batchRunner(fn WatchFunc) {
	t := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-t.C:
			batchMu.Lock()
			for filePath, lastEventTime := range batch {
				now := time.Now()
				if now.Sub(lastEventTime) > minBatchInterval {
					if err := fn(filePath); err != nil {
						log.Errorf("Failed to execute watch func: %v", err)
					}
					delete(batch, filePath)
				}
			}
			batchMu.Unlock()

		}
	}
}

func WatchDir(ctx context.Context, dir string, fn WatchFunc) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := watcher.Close(); err != nil {
			log.Errorf("Failed to close watcher cleanly: %v", err)
		}
	}()
	if err := watcher.Add(dir); err != nil {
		log.Fatalf("Failed to add watch dir: %v", dir)
	}
	go batchRunner(fn)
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-watcher.Events:
			if !ok {
				continue
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				batchMu.Lock()
				batch[event.Name] = time.Now()
				batchMu.Unlock()
				log.Debugf("Got write event: %s", event.Name)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				continue
			}
			log.Errorf("Watcher error: %v", err)
		}
	}
}

func init() {
	batch = make(map[string]time.Time)
	batchMu = &sync.RWMutex{}
}
