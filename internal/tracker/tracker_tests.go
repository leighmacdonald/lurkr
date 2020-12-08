// Package tracker implements a test suite for the tracker driver implementations to test against
// Note that this is intentionally named with _tests and not _test so its not automatically ran
package tracker

import (
	"github.com/leighmacdonald/lurkr/internal/config"
	"github.com/leighmacdonald/lurkr/internal/parser"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

const (
	// Auto replaced with real token when testing so we dont commit valid keys to git
	AuthKeyToken = "{{test.auth_key}}"
	PasskeyToken = "{{test.torrent_pass}}"
)

type TestData struct {
	Msg     string
	Res     *parser.Result
	Err     error
	Size    int64
	HashStr string
}

// TODO use transform.Format
func replaceAuthToken(driverName string, s *parser.Result) {
	if strings.Contains(s.LinkDL, AuthKeyToken) {
		c, _ := config.Tracker(driverName)
		s.LinkDL = strings.Replace(s.LinkDL, AuthKeyToken, c.Auth, 1)
	}
	if strings.Contains(s.LinkDL, PasskeyToken) {
		c, _ := config.Tracker(driverName)
		s.LinkDL = strings.Replace(s.LinkDL, PasskeyToken, c.Passkey, 1)
	}
}

func TestTracker(t *testing.T, driver Driver, testCases []TestData) {
	for i, c := range testCases {
		res, err := driver.ParseMessage(c.Msg)
		if c.Err != nil {
			require.Equal(t, c.Err, err)
			continue
		}
		require.NotNil(t, res)
		replaceAuthToken(c.Res.Tracker, res)
		require.NoErrorf(t, err, "%s %d failed to parse", driver, i)
		require.Equal(t, c.Res.Tracker, res.Tracker)
		require.Equal(t, c.Res.LinkSite, res.LinkSite)
		require.Equal(t, c.Res.Name, res.Name)
		require.Equal(t, c.Res.Category, res.Category)
		require.Equal(t, c.Res.Formats, res.Formats)
		require.Equal(t, c.Res.Year, res.Year)
		require.Equal(t, c.Res.SubName, res.SubName)
		require.Equal(t, c.Res.Tags, res.Tags)
		require.Equal(t, c.Res.Season, res.Season)
		require.Equal(t, c.Res.Episode, res.Episode)
		if c.Res.Group != "" {
			require.Equal(t, c.Res.Group, res.Group)
		}
		f, err := driver.Download(res)
		require.NoError(t, err, "Failed to download file")
		require.NotNil(t, f, "Empty file result")
		require.NotEqual(t, "", driver.Name(), "Invalid driver name")

		if c.Size > 0 || c.HashStr != "" {
			mi, err := f.UnmarshalInfo()
			require.NoError(t, err, "Failed to unmarshall meta info")
			s := mi.TotalLength()
			ih := f.HashInfoBytes().HexString()
			require.Equal(t, c.HashStr, ih)
			require.Equal(t, c.Size, s)
		}
	}
}
