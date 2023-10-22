package pw_test

import (
	"scraper/browserhandler/driver/pw"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBrowserIsConnected(t *testing.T) {
	driver := &pw.PwDriver{}
	driver.Start()
	browser := driver.Browser
	require.True(t, browser.IsConnected())
	driver.Close()
}
