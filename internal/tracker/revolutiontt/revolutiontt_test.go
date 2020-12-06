package revolutiontt

import (
	"github.com/leighmacdonald/lurkr/internal/config"
	"github.com/leighmacdonald/lurkr/internal/parser"
	"github.com/leighmacdonald/lurkr/internal/tracker"
	"github.com/stretchr/testify/require"
	"testing"
)

var tests = []tracker.TestData{
	{
		Msg: "!new Ellen.DeGeneres.2020.12.01.Justin.Bieber.720p.HDTV.x264-60FPS | " +
			"TV/HDx264 | " +
			"https://revolutiontt.me/details.php?id=ePZ7mjh3&hit=1",
		Res: &parser.Result{
			Tracker:  "revolutiontt",
			Name:     "Ellen DeGeneres 2020 12 01 Justin Bieber",
			SubName:  "",
			LinkSite: "https://revolutiontt.me/details.php?id=ePZ7mjh3&hit=1",
			LinkDL:   "",
			Year:     2020,
			Tags:     []string{"electronic", "folk", "country", "ambient", "dub", "world.music", "downtempo"},
			Formats:  []string{"mp3", "320", "cd"},
		},
		Err: nil,
	},
}

func TestRevolutionTT(t *testing.T) {
	tkr, err := New(&config.TrackerConfig{})
	require.NoErrorf(t, err, "Invalid tracker configuration: %s", driverName)
	tracker.TestTracker(t, tkr, tests)
}
