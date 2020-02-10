package resolver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPlaylistLocatorPageURL(t *testing.T) {
	playlistLocatorURL, err := GetPlaylistLocatorPageURL("881")
	if assert.NoError(t, err) {
		t.Logf("Playlist Locator URL: %s", playlistLocatorURL)
	}
}

func TestFind(t *testing.T) {
	playlist, err := Find("881")

	if assert.NoError(t, err) {
		t.Logf("Playlist: %+v", playlist)
	}
}
