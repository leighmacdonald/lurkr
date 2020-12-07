package config

import (
	"github.com/leighmacdonald/golib"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestTransportConfig(t *testing.T) {
	m, err := TransportConfigFile("default")
	require.NoError(t, err)
	require.NotNil(t, m)
	m2, err2 := TransportConfigSFTP("invalid")
	require.Error(t, err2)
	require.Nil(t, m2)

}

func TestMain(m *testing.M) {
	if err := Read(golib.FindFile("lurkr.yml", "lurkr")); err != nil {
		log.Errorf("Failed to read config")
		os.Exit(1)
	}
	os.Exit(m.Run())
}
