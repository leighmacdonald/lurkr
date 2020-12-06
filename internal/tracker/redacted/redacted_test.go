package redacted

import (
	"github.com/leighmacdonald/lurkr/internal/config"
	"github.com/leighmacdonald/lurkr/internal/parser"
	"github.com/leighmacdonald/lurkr/internal/tracker"
	"github.com/stretchr/testify/require"
	"testing"
)

var tests = []tracker.TestData{
	{
		Msg: "Various Artists - Elysian Vibes 5 [2001] [Compilation] - MP3 / 320 / CD - " +
			"https://redacted.ch/torrents.php?id=1440357 / " +
			"https://redacted.ch/torrents.php?action=download&id=3063049 - " +
			"electronic, folk, country, ambient, dub, world.music, downtempo",
		Res: &parser.Result{
			Tracker:  "redacted",
			Name:     "Various Artists",
			SubName:  "Elysian Vibes 5",
			LinkSite: "https://redacted.ch/torrents.php?id=3063049",
			LinkDL:   "3063049",
			Year:     2001,
			Tags:     []string{"electronic", "folk", "country", "ambient", "dub", "world.music", "downtempo"},
			Formats:  []string{"mp3", "320", "cd"},
		},
		Err: nil,
	},
}

func TestRedacted(t *testing.T) {
	tkr, err := New(&config.TrackerConfig{})
	require.NoErrorf(t, err, "Invalid tracker configuration: %s", driverName)
	tracker.TestTracker(t, tkr, tests)
}
