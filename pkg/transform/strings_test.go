package transform

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFormat(t *testing.T) {
	r1, e1 := Format("Test {test_1}", FArgs{"test_1": "replaced"})
	require.NoError(t, e1)
	require.Equal(t, "Test replaced", r1)

	_, e2 := Format("Test {test_2} {missing}", FArgs{"test_2": "replaced"})
	require.Error(t, e2)

	r3, e3 := Format("Test {test_3}", FArgs{"test_3": "replaced", "extra": 1})
	require.NoError(t, e3)
	require.Equal(t, "Test replaced", r3)
}
