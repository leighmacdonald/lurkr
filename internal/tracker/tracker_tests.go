// Package tracker implements a test suite for the tracker driver implementations to test against
// Note that this is intentionally named with _tests and not _test so its not automatically ran
package tracker

import (
	"github.com/leighmacdonald/lurkr/internal/parser"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestData struct {
	Msg string
	Res *parser.Result
	Err error
}

func TestTracker(t *testing.T, driver Driver, testCases []TestData) {
	for i, c := range testCases {
		res, err := driver.ParseMessage(c.Msg)
		if c.Err != nil {
			require.Equal(t, c.Err, err)
			continue
		}
		require.NoErrorf(t, err, "%s %d failed to parse", driver, i)
		require.Equal(t, c.Res, res)
		f, err := driver.Download(res)
		require.NoError(t, err, "Failed to download file")
		require.NotNil(t, f, "Empty file result")
	}
}
