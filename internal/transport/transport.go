package transport

import (
	"github.com/anacrolix/torrent/metainfo"
	"github.com/leighmacdonald/lurkr/internal/tracker"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

var (
	ErrInvalidTransport = errors.New("Invalid transport name")
)

type Transport interface {
	Send(file io.Reader, path string) error
}

// FetchTorrent provides a general purpose way of fetching a torrent file using
// the tracker specific http.Client. This
func FetchTorrent(client *http.Client, url string) (*metainfo.MetaInfo, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err2 := client.Do(req)
	if err2 != nil {
		return nil, err2
	}
	mi, err3 := metainfo.Load(resp.Body)
	if err3 != nil {
		return nil, tracker.ErrDownload
	}
	defer func() {
		if err4 := resp.Body.Close(); err4 != nil {
			log.Errorf("Failed to closed response body: %s", err4)
		}
	}()
	return mi, nil
}
