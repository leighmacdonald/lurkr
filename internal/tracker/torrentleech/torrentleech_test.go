package torrentleech

import (
	"fmt"
	"github.com/leighmacdonald/golib"
	"github.com/leighmacdonald/lurkr/internal/config"
	"github.com/leighmacdonald/lurkr/internal/parser"
	"github.com/leighmacdonald/lurkr/internal/tracker"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var tests = []tracker.TestData{
	{
		Msg: "New Torrent Announcement: <TV :: Episodes HD>  " +
			"Name:'American Monster S06E01 Everyones Favorite Uncle 1080p WEB h264-KOMPOST' " +
			"uploaded by 'Anonymous' -  https://www.torrentleech.org/torrent/1745788",
		Res: &parser.Result{
			Tracker:  driverName,
			Name:     "American Monster",
			SubName:  "",
			Group:    "KOMPOST",
			LinkSite: "https://www.torrentleech.org/torrent/1745788",
			LinkDL:   fmt.Sprintf("https://www.torrentleech.org/rss/download/1745795/%s/American.Monster.S06E01.Everyones.Favorite.Uncle.1080p.WEB.h264-KOMPOST.torrent", tracker.AuthKeyToken),
			Year:     0,
			Season:   6,
			Episode:  1,
			Category: parser.TV,
			Tags:     []string{"Episodes HD"},
			Formats:  []string{},
		},
		Err: nil,
	},
}

func TestTorrentLeech(t *testing.T) {
	tkr, err := New(config.Tracker(driverName))
	require.NoErrorf(t, err, "Invalid tracker configuration: %s", driverName)
	tracker.TestTracker(t, tkr, tests)
}

func TestMain(m *testing.M) {
	if err := config.Read(golib.FindFile("lurkr.yml", "lurkr")); err != nil {
		os.Exit(0)
	}
	os.Exit(m.Run())
}
