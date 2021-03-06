package revolutiontt

import (
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
			Category: parser.TV,
			Tags:     []string{"hdx264"},
			Formats:  []string{},
		},
		Err:     nil,
		Size:    1665820034,
		HashStr: "ffaa0efa86220b5065c3b5b9887a3d9764e07215",
	},
}

func TestRevolutionTT(t *testing.T) {
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
