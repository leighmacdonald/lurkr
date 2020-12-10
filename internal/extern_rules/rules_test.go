package extern_rules

import (
	"github.com/leighmacdonald/golib"
	"github.com/stretchr/testify/require"
	"path"
	"testing"
)

func TestImportIrssiRules(t *testing.T) {
	rootDir := golib.FindFile(path.Join("extern", "autodl-trackers", "trackers"), "lurkr")
	rules, err := ImportIrssiRules(rootDir)
	require.NoError(t, err, "Failed to import rules")
	require.Equal(t, 112, len(rules), "Invalid rule count")
	for _, rule := range rules {
		require.True(t, len(rule.ExampleLines) > 0, "Missing example: %s", rule.LongName)
	}
}
