package broadcasthenet

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
		Msg: "Ancient Aliens | S16E03 | Episode | 2020 | MKV | H.264 | WEB-DL | SD | Yes | Yes | " +
			"1390732 | Anonymous | English | Ancient.Aliens.S16E03.WEB.h264-BAE",
		Res: &parser.Result{
			Tracker:  driverName,
			Name:     "Ancient Aliens",
			SubName:  "",
			Group:    "BAE",
			LinkSite: "https://broadcasthe.net/torrents.php?torrentid=1390732",
			LinkDL: fmt.Sprintf(
				"https://broadcasthe.net/torrents.php?action=download&id=1390732&authkey=%s&torrent_pass=%s",
				tracker.AuthKeyToken, tracker.PasskeyToken),
			Year:     2020,
			Season:   16,
			Episode:  3,
			Category: parser.TV,
			Tags:     []string{},
			Formats:  []string{},
		},
		Err: nil,
	},
}

func TestBroadcasthenet(t *testing.T) {
	cfg, cfgErr := config.Tracker(driverName)
	require.NoErrorf(t, cfgErr, "Invalid tracker name: %s", driverName)
	tkr, err := New(cfg)
	require.NoErrorf(t, err, "Invalid tracker configuration: %s", driverName)
	tracker.TestTracker(t, tkr, tests)
}

func TestMain(m *testing.M) {
	if err := config.Read(golib.FindFile("lurkr.yml", "lurkr")); err != nil {
		os.Exit(0)
	}
	os.Exit(m.Run())
}
